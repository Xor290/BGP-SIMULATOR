package bgp

import (
	"bgp-manager/models"
	"fmt"
	"os"
	"os/exec"

	"gopkg.in/yaml.v3"
)

func RunAnsiblePlaybook(localAS *models.AutonomousSystem, peers []models.Peer) error {
	vars := models.PlaybookVars{
		LocalASN: localAS.ASN,
		RouterID: localAS.RouterID,
		Hostname: fmt.Sprintf("frr-as%d", localAS.ASN),
		Peers:    make([]models.PeerVars, 0, len(peers)),
	}

	for _, p := range peers {
		vars.Peers = append(vars.Peers, models.PeerVars{
			PeerIP:      p.PeerIP,
			RemoteASN:   p.RemoteASN,
			Description: p.Description,
			Password:    p.Password,
			Enabled:     p.Enabled,
		})
	}

	data, err := yaml.Marshal(vars)
	if err != nil {
		return fmt.Errorf("marshal vars: %w", err)
	}

	tmpFile, err := os.CreateTemp("", fmt.Sprintf("bgp-vars-as%d-*.yml", localAS.ASN))
	if err != nil {
		return fmt.Errorf("create temp file: %w", err)
	}
	defer os.Remove(tmpFile.Name())

	if _, err := tmpFile.Write(data); err != nil {
		return fmt.Errorf("write vars file: %w", err)
	}
	tmpFile.Close()

	out, err := exec.Command(
		"ansible-playbook",
		"ansible/apply_peer.yml",
		"-i", "ansible/inventory/hosts.yml",
		"-e", "@"+tmpFile.Name(),
	).CombinedOutput()
	if err != nil {
		return fmt.Errorf("ansible-playbook: %s: %w", out, err)
	}

	return nil
}
