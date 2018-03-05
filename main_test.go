package main

import (
	"database/sql"
	"os/exec"
	"testing"

	_ "github.com/lib/pq"
)

// Test multiple versions of Postgres
const (
	version94 = "9.4"
	version95 = "9.5"
	version96 = "9.6"
	version10 = "10"
)

const (
	command_11_4     = "cd /home/test/ && ./Setup.sh localhost m"
	command_11_5     = "cd /home/test/ && ./Setup.sh -h localhost -d m -s test -i 123"
	command_trunk    = "cd /home/test/ && ./Setup.sh -h localhost -d m -s test -i 123"
	sourcePath_11_4  = "/Users/pawel/Work/repo/SourceCode/D_11_4/Source/Scripts/test/"
	sourcePath_11_5  = "/Users/pawel/Work/repo/SourceCode/D_11_5/Source/Scripts/test/"
	sourcePath_trunk = "/Users/pawel/Work/repo/SourceCode/trunk/Source/Scripts/test/"
)

func dbSchemaTest_11_4(t *testing.T, conn *sql.DB, db *PostgresDB) {
	out, err := exec.Command("docker", "cp", sourcePath_11_4, db.cid+":/home/").CombinedOutput()
	if err != nil {
		t.Error(string(out), err)
	}
	com, err := exec.Command("docker", "container", "exec", db.cid, "bash", "-c", command_11_4).CombinedOutput()
	if err != nil {
		t.Error(string(com), err)
	}
	//fmt.Println(string(com))
}

func dbSchemaTest_11_5(t *testing.T, conn *sql.DB, db *PostgresDB) {
	out, err := exec.Command("docker", "cp", sourcePath_11_5, db.cid+":/home/").CombinedOutput()
	if err != nil {
		t.Error(string(out), err)
	}
	com, err := exec.Command("docker", "container", "exec", db.cid, "bash", "-c", command_11_5).CombinedOutput()
	if err != nil {
		t.Error(string(com), err)
	}
	//fmt.Println(string(com))
}

func dbSchemaTest_trunk(t *testing.T, conn *sql.DB, db *PostgresDB) {
	out, err := exec.Command("docker", "cp", sourcePath_trunk, db.cid+":/home/").CombinedOutput()
	if err != nil {
		t.Error(string(out), err)
	}
	com, err := exec.Command("docker", "container", "exec", db.cid, "bash", "-c", command_trunk).CombinedOutput()
	if err != nil {
		t.Error(string(com), err)
	}
	//fmt.Println(string(com))
}

// The actual Go tests just immediately call RunDBTest
//func TestCreateTable94(t *testing.T)  { RunDBTest(t, version94, testCreateTable) }

// 11.4
func TestPostgreSQL94_11_4(t *testing.T) { RunDBTest(t, version94, dbSchemaTest_11_4) }
func TestPostgreSQL95_11_4(t *testing.T) { RunDBTest(t, version95, dbSchemaTest_11_4) }
func TestPostgreSQL96_11_4(t *testing.T) { RunDBTest(t, version96, dbSchemaTest_11_4) }
func TestPostgreSQL10_11_4(t *testing.T) { RunDBTest(t, version10, dbSchemaTest_11_4) }

// 11.5
func TestPostgreSQL94_11_5(t *testing.T) { RunDBTest(t, version94, dbSchemaTest_11_5) }
func TestPostgreSQL95_11_5(t *testing.T) { RunDBTest(t, version95, dbSchemaTest_11_5) }
func TestPostgreSQL96_11_5(t *testing.T) { RunDBTest(t, version96, dbSchemaTest_11_5) }
func TestPostgreSQL10_11_5(t *testing.T) { RunDBTest(t, version10, dbSchemaTest_11_5) }

// trunk
func TestPostgreSQL94_trunk(t *testing.T) { RunDBTest(t, version94, dbSchemaTest_trunk) }
func TestPostgreSQL95_trunk(t *testing.T) { RunDBTest(t, version95, dbSchemaTest_trunk) }
func TestPostgreSQL96_trunk(t *testing.T) { RunDBTest(t, version96, dbSchemaTest_trunk) }
func TestPostgreSQL10(t *testing.T)       { RunDBTest(t, version10, dbSchemaTest_11_4) }
