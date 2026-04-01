# BGP Simulator

Simulateur de routage BGP reproduisant les mécanismes des grands opérateurs réseau.
Combine une API Go, des containers FRRouting, Ansible pour la configuration dynamique, et un frontend React.

---

## Sommaire

1. [Architecture générale](#architecture-générale)
2. [Le protocole BGP](#le-protocole-bgp)
3. [Autonomous Systems](#autonomous-systems)
4. [BGP inter-AS](#bgp-inter-as)
5. [Mitigation DDoS via BGP](#mitigation-ddos-via-bgp)
6. [API — Endpoints](#api--endpoints)
7. [Flux de configuration](#flux-de-configuration)
8. [Lancer le projet](#lancer-le-projet)

---

## Architecture générale

```mermaid
graph TD
    FE["Frontend React\nlocalhost:5173"]
    API["API Go · Gin\nlocalhost:8080"]
    DB["PostgreSQL\nautonomous_systems\npeers · prefix_since_as\nbgp_sessions"]
    ANS["Ansible\napply_peer.yml"]
    FRR1["frr-as65001\nFRRouting 9.1\nAS 65001 · 10.0.0.11"]
    FRR2["frr-as65002\nFRRouting 9.1\nAS 65002 · 10.0.0.12"]

    FE -->|HTTP + JWT| API
    API -->|GORM| DB
    API -->|génère vars + lance| ANS
    ANS -->|Docker connection\nfrr.conf| FRR1
    ANS -->|Docker connection\nfrr.conf| FRR2
    FRR1 <-->|eBGP TCP 179\nbgp-fabric 10.0.0.0/24| FRR2
```

---

## Le protocole BGP

**BGP (Border Gateway Protocol)** est le protocole de routage qui fait fonctionner Internet.
Défini par la [RFC 4271](https://datatracker.ietf.org/doc/html/rfc4271), il opère sur **TCP port 179**.

Contrairement aux protocoles IGP (OSPF, IS-IS) qui cherchent le chemin le plus court,
BGP est un **protocole à vecteur de chemin** — il choisit la route selon des **politiques** (policies).

### Machine à états (FSM — RFC 4271 §8)

```mermaid
stateDiagram-v2
    [*] --> Idle
    Idle --> Connect : Start event
    Connect --> OpenSent : TCP established
    Connect --> Active : TCP failed
    Active --> OpenSent : TCP established
    OpenSent --> OpenConfirm : OPEN received & valid
    OpenConfirm --> Established : KEEPALIVE received
    Established --> Idle : Error / Hold timer expired
    OpenSent --> Idle : Error
    OpenConfirm --> Idle : Error
```

| État          | Description                                                              |
|---------------|--------------------------------------------------------------------------|
| `Idle`        | BGP inactif, en attente d'un événement de démarrage                     |
| `Connect`     | Tentative de connexion TCP vers le voisin (port 179)                    |
| `Active`      | TCP échoué, BGP réessaie activement                                     |
| `OpenSent`    | TCP établi, message OPEN envoyé — attente de l'OPEN du voisin           |
| `OpenConfirm` | OPEN reçu et validé — attente du KEEPALIVE final                        |
| `Established` | Session active, échange de UPDATE et KEEPALIVE                          |

### Les 4 messages BGP

| Message        | Rôle                                                                    |
|----------------|-------------------------------------------------------------------------|
| `OPEN`         | Ouvre la session : ASN local, Router ID, hold time, capacités           |
| `UPDATE`       | Annonce ou retire des préfixes (NLRI) avec leurs attributs de chemin   |
| `KEEPALIVE`    | Maintien de la session (~60s), confirme aussi l'OPEN                   |
| `NOTIFICATION` | Signale une erreur fatale, ferme immédiatement la session              |

### Attributs de chemin

| Attribut     | Type        | Rôle                                                              |
|--------------|-------------|-------------------------------------------------------------------|
| `AS_PATH`    | Obligatoire | Liste des AS traversés — évite les boucles                       |
| `NEXT_HOP`   | Obligatoire | IP du prochain saut pour atteindre le préfixe                    |
| `LOCAL_PREF` | Optionnel   | Préférence locale (plus élevé = préféré), propagé en iBGP        |
| `MED`        | Optionnel   | Indique au voisin le chemin entrant préféré                      |
| `ORIGIN`     | Obligatoire | Origine : IGP (`i`), EGP (`e`), Incomplete (`?`)                 |

### Sélection de route (Decision Process)

BGP choisit la meilleure route dans cet ordre :

1. `LOCAL_PREF` le plus élevé
2. `AS_PATH` le plus court
3. `ORIGIN` le plus bas (IGP < EGP < Incomplete)
4. `MED` le plus faible
5. eBGP préféré sur iBGP
6. IGP metric la plus faible vers le NEXT_HOP
7. Router ID le plus bas (tie-breaker)

---

## Autonomous Systems

Un **Autonomous System** est un ensemble de réseaux IP sous une même administration,
identifié par un **numéro unique (ASN)**.

| Plage                   | Usage                                    |
|-------------------------|------------------------------------------|
| 1 – 64495               | ASN publics 16 bits (IANA/RIR)          |
| 64512 – 65535           | ASN privés 16 bits (labo/simulation)    |
| 65536 – 4 294 967 295   | ASN 32 bits (RFC 4893)                  |

Exemples réels : Cloudflare = **AS13335** · Google = **AS15169** · Hurricane Electric = **AS6939**

### Modèle en base

```mermaid
erDiagram
    AutonomousSystem {
        uint ID PK
        uint32 ASN
        string Name
        string RouterID
        string Description
    }
    Peer {
        uint ID PK
        uint LocalASID FK
        uint32 RemoteASN
        string PeerIP
        string Description
        bool Enabled
    }
    PrefixSinceAS {
        uint ID PK
        uint ASID FK
        string Prefix
        string NextHop
        uint LocalPref
        uint MED
        bool Active
    }
    BGPSession {
        uint ID PK
        uint PeerID FK
        string State
        string RemoteIP
        int MsgRcvd
        int MsgSent
    }

    AutonomousSystem ||--o{ Peer : "has"
    AutonomousSystem ||--o{ PrefixSinceAS : "announces"
    Peer ||--o{ BGPSession : "has"
```

---

## BGP inter-AS

### eBGP vs iBGP

| Type  | Contexte                          | NEXT_HOP          | LOCAL_PREF |
|-------|-----------------------------------|-------------------|------------|
| eBGP  | Entre deux AS différents          | Mis à jour        | Non propagé |
| iBGP  | Au sein du même AS                | Non modifié       | Propagé    |

Ce projet simule uniquement de l'**eBGP**.

### Échange de messages lors d'une session

```mermaid
sequenceDiagram
    participant A as frr-as65001 (AS 65001)
    participant B as frr-as65002 (AS 65002)

    A->>B: TCP SYN → port 179
    B->>A: TCP SYN-ACK
    A->>B: OPEN (ASN=65001, RouterID=10.0.0.11)
    B->>A: OPEN (ASN=65002, RouterID=10.0.0.12)
    A->>B: KEEPALIVE
    B->>A: KEEPALIVE
    Note over A,B: État : Established
    A->>B: UPDATE (192.168.1.0/24, NEXT_HOP=10.0.0.11, AS_PATH=[65001])
    B->>A: UPDATE (10.10.0.0/24, NEXT_HOP=10.0.0.12, AS_PATH=[65002])
    loop Toutes les ~60s
        A->>B: KEEPALIVE
        B->>A: KEEPALIVE
    end
```

### Configuration FRR générée

Pour qu'AS65001 annonce `192.168.1.0/24` à AS65002 :

```
router bgp 65001
 bgp router-id 10.0.0.11
 no bgp ebgp-requires-policy
 neighbor 10.0.0.12 remote-as 65002
 !
 address-family ipv4 unicast
  network 192.168.1.0/24
  neighbor 10.0.0.12 activate
 exit-address-family
```

> `no bgp ebgp-requires-policy` : FRR 9.x bloque par défaut l'échange de préfixes
> sans route-map explicite (`(Policy)` dans `show bgp summary`). Cette directive lève cette restriction.

---

## Mitigation DDoS via BGP

Quand un serveur est sous attaque DDoS, BGP permet de **rediriger le trafic malveillant**
vers une infrastructure de nettoyage (scrubbing center) avant de renvoyer le trafic légitime
vers le serveur cible via un **scrubbing center**.

### Flux normal (sans attaque)

```mermaid
flowchart LR
    Internet["Internet\n(trafic légitime)"]
    RT["Routeur de bordure\nAS65000"]
    SRV["Serveur cible\n1.2.3.4"]

    Internet -->|"BGP annonce\n1.2.3.4/32"| RT
    RT --> SRV
```

### Flux sous attaque — BGP Blackhole

Le cas le plus simple : l'IP attaquée est **blackholée** (tout le trafic vers elle est jeté).
Utilisé quand l'attaque est si volumineuse qu'elle sature les liens.

```mermaid
flowchart LR
    Internet["Internet\nDDoS 500 Gbps"]
    RT["Routeur de bordure\nAS65000"]
    BH["Blackhole\n/dev/null"]
    SRV["Serveur cible\n1.2.3.4\n❌ inaccessible"]

    Internet -->|"UPDATE BGP\ncommunity 16276:666\n1.2.3.4/32 → Null0"| RT
    RT --> BH
    RT -. trafic stoppé .-> SRV

    style BH fill:#7f1d1d,color:#fca5a5
    style SRV fill:#1c1917,color:#78716c
```

> La community BGP `65000:666` est un signal envoyé aux upstreams.
> Dès réception, chaque opérateur transit jette le trafic vers cette IP avant qu'il n'atteigne le réseau.

### Flux sous attaque — BGP Rerouting vers VAC

Solution plus fine : le trafic est **dévié vers le scrubbing center**, nettoyé, puis réinjecté.

```mermaid
flowchart TD
    Internet["Internet\nDDoS + trafic légitime"]
    RT["Routeur de bordure\nAS65000"]
    VAC["VAC · Scrubbing Center\nAnalyse + filtrage\n1.2.3.4/32 via tunnel"]
    SRV["Serveur cible\n1.2.3.4\n✅ reçoit trafic propre"]
    DETECT["Système de détection\nseuil volumétrique"]

    Internet --> RT
    DETECT -->|"Attaque détectée\nBGP UPDATE\n1.2.3.4/32 → VAC"| RT
    RT -->|"Tout le trafic\nredirigé"| VAC
    VAC -->|"Trafic légitime\nré-injecté via GRE/MPLS"| SRV
    VAC -->|"Trafic malveillant\ndroppé"| DROPPED["❌ Dropped"]

    style VAC fill:#14532d,color:#86efac
    style SRV fill:#14532d,color:#86efac
    style DROPPED fill:#7f1d1d,color:#fca5a5
```

### Mécanisme BGP utilisé

```mermaid
sequenceDiagram
    participant DET as Détection DDoS
    participant RT  as Routeur bordure (AS65000)
    participant UP  as Upstream (Telia, NTT...)
    participant VAC as Scrubbing Center

    Note over RT,UP: État normal — 1.2.3.4/32 annoncé normalement
    DET->>RT: Seuil dépassé sur 1.2.3.4
    RT->>RT: Modifie next-hop de 1.2.3.4/32\nvers IP du VAC
    RT->>UP: BGP UPDATE\n1.2.3.4/32 community 16276:9999
    UP->>VAC: Trafic dévié vers scrubbing center
    VAC->>RT: Trafic propre réinjecté (tunnel GRE)
    Note over RT,UP: Fin d'attaque
    RT->>UP: BGP UPDATE\n1.2.3.4/32 next-hop normal (withdraw mitigation)
```

### Les deux techniques comparées

| Technique         | BGP Blackhole                        | BGP Rerouting (VAC)                    |
|-------------------|--------------------------------------|----------------------------------------|
| Trafic légitime   | ❌ Jeté avec le malveillant           | ✅ Nettoyé et réacheminé               |
| Vitesse           | Immédiate (quelques secondes)        | Rapide mais avec analyse (~30s)        |
| Cas d'usage       | Attaque volumétrique extrême         | Attaque standard, service à maintenir  |
| Propagation       | Chez tous les upstreams via community | Interne à l'AS ou vers upstreams proches |

---

## API — Endpoints

Tous les endpoints sauf `/health` et `/auth/*` requièrent :
```
Authorization: Bearer <access_token>
```

### Auth

#### `POST /api/v1/auth/register`
```json
{
  "username": "admin",
  "password": "motdepasse",
  "nom": "Dupont",
  "prenom": "Jean",
  "telephone": "+33600000000"
}
```
Réponse `201` : `{ "message": "Client enregistré" }`

#### `POST /api/v1/auth/login`
```json
{ "username": "admin", "password": "motdepasse" }
```
Réponse `200` :
```json
{ "access_token": "eyJhbGci...", "token_type": "Bearer", "expires_in": 3600 }
```

---

### Autonomous Systems

#### `POST /api/v1/as/`
```json
{ "asn": 65001, "name": "AS-65001", "router_id": "10.0.0.11" }
```
Réponse `201` : objet `AutonomousSystem`.

---

### Peers

| Méthode   | Route                          | Description                          |
|-----------|-------------------------------|--------------------------------------|
| `GET`     | `/api/v1/peers/all`           | Lister tous les peers                |
| `POST`    | `/api/v1/peers/create`        | Créer un peer → déclenche Ansible   |
| `GET`     | `/api/v1/peers/:id`           | Détail d'un peer                     |
| `DELETE`  | `/api/v1/peers/:id`           | Supprimer un peer                    |
| `GET`     | `/api/v1/peers/:id/sessions`  | Sessions BGP du peer                 |

#### Body `POST /api/v1/peers/create`
```json
{
  "local_as_id": 1,
  "remote_asn": 65002,
  "peer_ip": "10.0.0.12",
  "description": "AS65001 -> AS65002",
  "enabled": true
}
```

---

### Préfixes

#### `POST /api/v1/bgp/create/prefix`
Stocke le préfixe en base et déclenche Ansible pour l'annoncer dans FRR.

```json
{
  "prefix": "192.168.1.0/24",
  "asn": 65001,
  "next_hop": "10.0.0.11",
  "local_pref": 100
}
```
Réponse `201` : objet `PrefixSinceAS`.

---

### Health

#### `GET /api/v1/health`
```json
{ "status": "ok" }
```

---

## Flux de configuration

```mermaid
flowchart TD
    REQ["Requête HTTP\nPOST /peers/create\nou /bgp/create/prefix"]
    VALID["Validation binding\n+ écriture PostgreSQL"]
    APPLY["ApplyASConfig\nasID"]
    QUERY["GetPeersByASID\nGetPrefixesByASID"]
    VARS["GenerateVarsFile\nYAML temporaire"]
    PLAYBOOK["ansible-playbook\napply_peer.yml"]
    TEMPLATE["Template Jinja2\nfrr_bgp.j2"]
    CONF["/etc/frr/frr.conf\nmis à jour"]
    RELOAD["vtysh write\n+ reload FRR"]
    SESSION["Session BGP\nmise à jour"]

    REQ --> VALID --> APPLY --> QUERY --> VARS --> PLAYBOOK --> TEMPLATE --> CONF --> RELOAD --> SESSION
```

---

## Lancer le projet

### Prérequis
- Docker + Docker Compose
- Go 1.22+
- Node 18+

### Stack complète

```bash
docker compose up -d --build
```

### Frontend

```bash
cd frontend && npm install && npm run dev
# http://localhost:5173
```

### Debug FRR

```bash
# Sessions BGP
docker exec frr-as65001 vtysh --vty_socket /var/run/frr -c "show bgp summary"

# Table de préfixes
docker exec frr-as65001 vtysh --vty_socket /var/run/frr -c "show bgp ipv4 unicast"

# Config active
docker exec frr-as65001 cat /etc/frr/frr.conf
```
