package bgp

import (
	"bgp-manager/models"
	tmpl "bgp-manager/template"
	"fmt"
	"os"
	"os/exec"
)

func RunAnsiblePlaybook(localAS *models.AutonomousSystem, peers []models.Peer) error {
	varsFile, err := tmpl.GenerateVarsFile(localAS, peers)
	if err != nil {
		return err
	}
	defer os.Remove(varsFile)

	cmd := exec.Command(
		"ansible-playbook",
		"ansible/apply_peer.yml",
		"-i", "ansible/inventory/hosts.yml",
		"-e", "@"+varsFile,
	)
	cmd.Env = append(os.Environ(), "ANSIBLE_CONFIG=ansible/ansible.cfg")
	out, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("ansible-playbook: %s: %w", out, err)
	}

	return nil
}
