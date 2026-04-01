echo "BGP summary:"
docker exec frr-as65001 vtysh --vty_socket /var/run/frr -c "show bgp summary" 2>&1 | grep -v "vtysh.conf"
echo ""
docker exec frr-as65002 vtysh --vty_socket /var/run/frr -c "show bgp summary" 2>&1 | grep -v "vtysh.conf"

echo ""
echo "BGP routes annonce:"
docker exec frr-as65001 vtysh --vty_socket /var/run/frr -c "show bgp ipv4 unicast" 2>&1 | grep -v "vtysh.conf"
echo ""
docker exec frr-as65002 vtysh --vty_socket /var/run/frr -c "show bgp ipv4 unicast" 2>&1 | grep -v "vtysh.conf"

echo ""
echo "BGP routes:"
docker exec frr-as65001 vtysh --vty_socket /var/run/frr -c "show ip bgp" 2>&1 | grep -v "vtysh.conf"
echo ""
docker exec frr-as65002 vtysh --vty_socket /var/run/frr -c "show ip bgp" 2>&1 | grep -v "vtysh.conf"

echo ""
echo "BGP neighbors:"
docker exec frr-as65001 vtysh --vty_socket /var/run/frr -c "show ip bgp neighbors" 2>&1 | grep -v "vtysh.conf"
echo ""
docker exec frr-as65002 vtysh --vty_socket /var/run/frr -c "show ip bgp neighbors" 2>&1 | grep -v "vtysh.conf"

echo ""
echo "Table in base:"
docker exec bgp-simulator-postgres-1 psql -U bgp -d bgp_manager -c "SELECT id, prefix, as_id, next_hop, active FROM prefix_since_a_s;"
docker exec bgp-simulator-postgres-1 psql -U bgp -d bgp_manager -c "SELECT id, asn, name, router_id FROM autonomous_systems;"
docker exec bgp-simulator-postgres-1 psql -U bgp -d bgp_manager -c "SELECT id, peer_ip, remote_asn, enabled FROM peers;"

echo ""
echo "Test ping:"
docker exec frr-as65001 ping -c 2 10.0.0.12
docker exec frr-as65002 ping -c 2 10.0.0.11
