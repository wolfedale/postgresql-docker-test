package main

import (
	"encoding/json"
	"fmt"
	"os/exec"
	"strings"
)

type PostgresDB struct {
	// Of form 'host:port'
	Host string

	// The container id assigned by docker
	cid string
}

type PostgresConfig struct {
	Password string
	Username string // defaults to "postgres"
	Database string // defaults to "username"
	Version  string // defaults to "latest"
}

func NewPostgresDB(c PostgresConfig) (*PostgresDB, error) {
	img := "postgres:latest"
	if c.Version == "" {
		img = "postgres:" + c.Version
	}

	// docker's run command has the nasty habbit of pulling images if you don't have them.
	// Warn user they need to pull the image, don't automatically pull for them.
	if exec.Command("docker", "inspect", img).Run() != nil {
		return nil, fmt.Errorf("db requires docker image %s, please pull or specify a different version", img)
	}

	// Running on port 0 instructs the operating system to pick an available port.
	dockerArgs := []string{"run", "-d", "-p", "127.0.0.1:0:5432"}
	envvars := map[string]string{
		"POSTGRES_PASSWORD": c.Password,
		"POSTGRES_USER":     c.Username,
		"POSTGRES_DB":       c.Database,
	}
	for key, val := range envvars {
		if val != "" {
			dockerArgs = append(dockerArgs, "-e", key+"="+val)
		}
	}
	dockerArgs = append(dockerArgs, img)

	// Start the docker container.
	out, err := exec.Command("docker", dockerArgs...).CombinedOutput()
	if err != nil {
		return nil, fmt.Errorf("docker run: %v: %s", err, out)
	}

	cid := strings.TrimSpace(string(out))
	db := &PostgresDB{cid: cid}

	db.Host, err = portMapping(cid, "5432/tcp")
	if err != nil {
		db.Close()
		return nil, err
	}
	return db, nil
}

func portMapping(cid, containerPort string) (hostAddr string, err error) {
	out, err := exec.Command("docker", "inspect", cid).CombinedOutput()
	if err != nil {
		return "", fmt.Errorf("docker inspect: %v: %s", err, out)
	}

	// anonymous struct for unmarshalling JSON into
	var inspectResp []struct {
		NetworkSettings struct {
			Ports map[string][]struct {
				HostIp   string
				HostPort string
			}
		}
	}
	if err := json.Unmarshal(out, &inspectResp); err != nil {
		return "", fmt.Errorf("decoding docker inspect result failed: %v: %s", err, out)
	}
	if len(inspectResp) != 1 {
		return "", fmt.Errorf("expected one inspect result, got %d", len(inspectResp))
	}
	ports := inspectResp[0].NetworkSettings.Ports[containerPort]
	if len(ports) != 1 {
		return "", fmt.Errorf("expected one port mapping, got %d", len(ports))
	}
	return ports[0].HostIp + ":" + ports[0].HostPort, nil
}

// Close removes the container running the postgres database.
func (db *PostgresDB) Close() error {
	out, err := exec.Command("docker", "rm", "-f", db.cid).CombinedOutput()
	if err != nil {
		return fmt.Errorf("docker rm: %v: %s", err, out)
	}
	return nil
}
