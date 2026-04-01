package bgp

import (
	"fmt"
	"os/exec"
)

func RunVtyshCommand(asn uint32, cmd string) error {
	containerName := fmt.Sprintf("frr-as%d", asn)
	out, err := exec.Command("docker", "exec", containerName, "vtysh", "-c", cmd).CombinedOutput()
	if err != nil {
		return fmt.Errorf("vtysh error: %s: %w", out, err)
	}
	return nil
}
