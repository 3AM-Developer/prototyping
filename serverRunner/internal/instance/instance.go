package instance

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"os"
	"os/exec"
	"path"

	"github.com/gorcon/rcon"
)

type InstanceClient interface {
	Stop()
	Start()
}

type Instance struct {
	Id          int    `json:"-"`
	Name        string `json:"name"`
	Dir         string `json:"-"`
	JavaPath    string `json:"java-path"`
	StartScript string `json:"start-script"`
	RconPort    int    `json:"rcon-port"`
	RconPw      string `json:"rcon-pw"`
}

func (i *Instance) Copy() *Instance {
	copy := Instance{
		Id:          i.Id,
		Name:        i.Name,
		Dir:         i.Dir,
		JavaPath:    i.JavaPath,
		StartScript: i.StartScript,
		RconPort:    i.RconPort,
		RconPw:      i.RconPw,
	}

	return &copy
}

func (i *Instance) getFullScriptPath() string {
	if path.IsAbs(i.StartScript) {
		return i.StartScript
	}

	return path.Join(i.Dir, i.StartScript)
}

func (i *Instance) Start() error {
	// TODO(jaegyu): Cleanup
	// We have to interface slightly with the OS layer here.
	// I want to clean this up

	screenCmd := exec.Command("screen", "-S", i.Name, "-d", "-m",
		i.getFullScriptPath())

	screenCmd.Stdin = os.Stdin
	screenCmd.Stdout = os.Stdout
	screenCmd.Stderr = os.Stderr

	err := screenCmd.Start()
	if err != nil {
		return err
	}

	return nil
}

func (i *Instance) Stop() error {
	conn, err := rcon.Dial(
		fmt.Sprintf("localhost:%d", i.RconPort),
		i.RconPw,
	)
	if err != nil {
		return err
	}

	response, err := conn.Execute("stop")
	if err != nil {
		return err
	}

	log.Println(response)
	return nil
}

// VerifyInstance will ignore missing id fields
// as this is primaryily used to validate before writing to file
func (i *Instance) VerifyInstance() bool {
	return i.Name != "" &&
		i.Dir != "" &&
		i.JavaPath != "" &&
		i.StartScript != "" &&
		i.RconPw != ""
}

func (i *Instance) Write() error {
	data, err := json.Marshal(i)
	if err != nil {
		return err
	}

	outfile := fmt.Sprintf(i.Dir, "instance.json")
	err = os.WriteFile(outfile, data, 640)
	if err != nil {
		return err
	}

	return nil
}

// TODO(jaegyu): finish loadconfig function
func LoadJson(i *Instance, path string) error {
	file, err := os.Open(path)
	if err != nil {
		return err
	}

	data, err := io.ReadAll(file)
	if err != nil {
		return err
	}

	// copy the instance to prevent any weird sideeffects on err
	var newI = i.Copy()

	err = json.Unmarshal(data, &i)
	if err != nil {
		i = newI
		return err
	}

	return nil
}
