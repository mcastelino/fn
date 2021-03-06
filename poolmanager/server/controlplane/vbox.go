package controlplane

import (
	"bytes"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"os"
	"os/exec"
	"strings"

	idgen "github.com/fnproject/fn/api/id"
)

const vboxNamePrefix = "fn-vagrant"

var whichVBox *exec.Cmd

func init() {
	whichVBox = exec.Command("which", "vagrant")
}

type VirtualBoxCP struct {
	runnerMap   map[string][]*Runner
	vagrantPath string
}

func NewVirtualBoxCP(vagrantPath string) (*VirtualBoxCP, error) {
	runnerMap := make(map[string][]*Runner)
	if err := whichVBox.Run(); err != nil {
		return nil, err
	}
	return &VirtualBoxCP{
		runnerMap:   runnerMap,
		vagrantPath: vagrantPath,
	}, nil
}

func (v *VirtualBoxCP) provision() (*Runner, error) {
	//set up dir
	wd, err := os.Getwd()
	if err != nil {
		return nil, err
	}
	defer func() {
		os.Chdir(wd)
	}()

	node := newNodeName()
	nodeDir, err := ioutil.TempDir(wd, node)
	if err != nil {
		return nil, err
	}
	//copy vagrant file into there
	newVagrantFile := fmt.Sprintf("%s/%s", nodeDir, "Vagrantfile")
	err = copyFile(v.vagrantPath, newVagrantFile)
	if err != nil {
		return nil, err
	}

	err = os.Chdir(nodeDir)
	if err != nil {
		log.Println(err.Error())
		return nil, err
	}
	vboxProvision := exec.Command("vagrant", "up")
	err = vboxProvision.Run()
	if err != nil {
		log.Println(err.Error())
		return nil, err
	}
	//Get the broadcast addr and call it a day
	return getRunner(node)
}

//Gets the address that its broadcasting at
//VBoxManage guestproperty get "cp_default_1520116902053_77841" "/VirtualBox/GuestInfo/Net/1/V4/Broadcast"
func getRunner(node string) (*Runner, error) {
	//TODO make the vagrant file templated
	vmsCmd := exec.Command("VBoxManage", "list", "vms")
	var vmsOut bytes.Buffer
	vmsCmd.Stdout = &vmsOut
	err := vmsCmd.Run()
	if err != nil {
		return nil, err
	}
	vms := strings.Split(vmsOut.String(), "\n")
	var realNode string
	for _, candidate := range vms {
		if strings.Contains(candidate, node) {
			spl := strings.Split(candidate, " ")
			realNode = spl[0]
		}
	}
	//strip the quotes
	if strings.Contains(realNode, "\"") {
		realNode = realNode[1 : len(realNode)-1]
	}

	//guestproperty get "fn-vagrant-6ae28c23-445e-4b0b-a2cf-0102e66ec57a766389779_default_1520288274551_74039" /VirtualBox/GuestInfo/Net/1/V4/Broadcast
	args := []string{"guestproperty", "get", realNode, "/VirtualBox/GuestInfo/Net/1/V4/Broadcast"}
	broadCastAddrCmd := exec.Command("VBoxManage", args...)
	var out bytes.Buffer
	broadCastAddrCmd.Stdout = &out

	var stdErr bytes.Buffer
	broadCastAddrCmd.Stderr = &stdErr

	err = broadCastAddrCmd.Run()
	if err != nil {
		log.Println("error running", err.Error(), stdErr.String())
		return nil, err
	}
	addr := strings.Split(out.String(), ":")
	if len(addr) != 2 {
		return nil, fmt.Errorf("Unable to get address got:'%s' as output", out.String())
	}
	return &Runner{
		Id:      realNode,
		Address: addr[1],
	}, nil
}

func (v *VirtualBoxCP) GetLBGRunners(lgbID string) ([]*Runner, error) {
	runners, ok := v.runnerMap[lgbID]
	if !ok {
		return nil, errors.New("Not Found")
	}
	return runners, nil
}

func (v *VirtualBoxCP) ProvisionRunners(lgbID string, n int) (int, error) {
	runners := make([]*Runner, 0, n)
	for i := 0; i < n; i++ {
		runner, err := v.provision()
		runners = append(runners, runner)
		if err != nil {
			return 0, err
		}
	}
	v.runnerMap[lgbID] = runners
	return n, nil
}

func (v *VirtualBoxCP) RemoveRunner(lbgID string, id string) error {
	wd, err := os.Getwd()
	if err != nil {
		return err
	}
	defer func() {
		os.Chdir(wd)
	}()

	runners, ok := v.runnerMap[lbgID]
	if !ok {
		return errors.New("No lgbID with this name")
	}
	//look for it in the customers map
	found := false
	for _, r := range runners {
		if id == r.Id {
			found = true
			break
		}
	}
	if found == false {
		return errors.New("No VM by this ID")
	}
	//switch to the dir and remove it
	//vm name is fn-vagrant-7183faa4-7321-47e9-8fd9-4a0aa1ac818e497509110_default_1520299457972_92567 everything before the first _
	split := strings.Split(id, "_")
	dirName := split[0]
	err = os.Chdir(dirName)
	if err != nil {
		log.Println(err.Error())
		return err
	}
	destroyCmd := exec.Command("vagrant", "destroy", "-f")
	err = destroyCmd.Run()
	if err != nil {
		log.Println(err.Error())
		return err
	}
	// back to working dir and rm -rf ignore these erro
	err = os.Chdir(wd)
	if err != nil {
		return err
	}
	err = os.RemoveAll(dirName)
	if err != nil {
		return err
	}
	return nil
}

func newNodeName() string {
	id := idgen.New()
	return fmt.Sprintf("%s-%s", vboxNamePrefix, id.String())
}

//TODO move to a util folder if needed again
func copyFile(src string, dst string) error {
	// Open the source file for reading
	s, err := os.Open(src)
	if err != nil {
		return err
	}
	defer s.Close()

	// Open the destination file for writing
	d, err := os.Create(dst)
	if err != nil {
		return err
	}

	// Copy the contents of the source file into the destination file
	if _, err := io.Copy(d, s); err != nil {
		d.Close()
		return err
	}

	// Return any errors that result from closing the destination file
	// Will return nil if no errors occurred
	return d.Close()
}
