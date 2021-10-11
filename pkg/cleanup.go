package gont

import (
	"fmt"
	"io/ioutil"
	"math/rand"
	"os"
	"path"
	"path/filepath"
	"sort"

	"github.com/vishvananda/netns"
)

func GetNetworkNames() []string {
	names := []string{}

	nets, err := ioutil.ReadDir(varDir)
	if err != nil {
		return names
	}

	for _, net := range nets {
		if net.IsDir() {
			names = append(names, net.Name())
		}
	}

	sort.Strings(names)

	return names
}

func GetNodeNames(network string) []string {
	names := []string{}

	nodesDir := path.Join(varDir, network, "nodes")

	nets, err := ioutil.ReadDir(nodesDir)
	if err != nil {
		return names
	}

	for _, net := range nets {
		if net.IsDir() {
			names = append(names, net.Name())
		}
	}

	sort.Strings(names)

	return names
}

func GenerateNetworkName() string {
	existing := GetNetworkNames()

	for i := 0; i < 32; i++ {
		random := GetRandomName()

		index := sort.SearchStrings(existing, random)
		if index >= len(existing) || existing[index] != random {
			return random
		}
	}

	index := rand.Intn(len(Names))
	random := Names[index]

	return fmt.Sprintf("%s%d", random, rand.Intn(128)+1)
}

func CleanupAllNetworks() error {
	for _, name := range GetNetworkNames() {
		if err := CleanupNetwork(name); err != nil {
			return err
		}
	}

	return nil
}

func CleanupNetwork(name string) error {
	baseDir := filepath.Join(varDir, name)
	nodesDir := filepath.Join(baseDir, "nodes")

	fis, err := ioutil.ReadDir(nodesDir)
	if err != nil {
		return err
	}

	for _, fi := range fis {
		if !fi.IsDir() {
			continue
		}

		nodeName := fi.Name()
		netNsName := fmt.Sprintf("gont-%s-%s", name, nodeName)

		netns.DeleteNamed(netNsName)
	}

	if err := os.RemoveAll(baseDir); err != nil {
		return err
	}

	return nil
}
