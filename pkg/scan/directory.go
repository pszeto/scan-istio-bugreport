package scan

import (
	"fmt"
	"io/ioutil"
	"os"
	"strings"

	log "github.com/sirupsen/logrus"
)

func getUniqueValue(deployments []string) map[string]int {
	//Create a   dictionary of values for each element
	dict := make(map[string]int)
	for _, deployment := range deployments {
		dict[deployment] = dict[deployment] + 1
	}
	return dict
}

func scanDirForDeployments(proxiesDir string, nsDirs []string) ([]NamespaceInfo, error) {
	var namespaceInfo []NamespaceInfo
	for _, f := range nsDirs {
		// log.Infof("Scanning directory : %s", proxiesDir+"/"+f)
		files, err := ioutil.ReadDir(proxiesDir + "/" + f)
		if err != nil {
			log.Errorf("Failed to read sub directory, error: %w", err)
			return nil, err
		}
		var deploymenNames []string
		for _, file := range files {
			if file.IsDir() {
				dirName := file.Name()
				dirName = dirName[0:strings.LastIndex(dirName, "-")]
				dirName = dirName[0:strings.LastIndex(dirName, "-")]
				deploymenNames = append(deploymenNames, dirName)
			}
		}
		namespaceInfo = append(namespaceInfo, NamespaceInfo{
			Name:        f,
			Deployments: getUniqueValue(deploymenNames),
		})
		deploymenNames = nil
	}

	return namespaceInfo, nil
}

func ScanForNsAndDeployments(dirName string, generate bool) error {
	proxiesDir := dirName + "/proxies"
	_, err := os.Stat(proxiesDir)
	if err != nil {
		log.Errorf("Failed to open proxies directory, error: %w", err)
		return err
	}

	files, err := ioutil.ReadDir(proxiesDir)
	if err != nil {
		log.Errorf("Failed to read proxies directory, error: %w", err)
		return err
	}
	var nsDirs []string
	for _, file := range files {
		if file.IsDir() && file.Name() != "istio-system" && file.Name() != "istio-gateways" {
			nsDirs = append(nsDirs, file.Name())
		}
	}

	namespaceInfo, err := scanDirForDeployments(proxiesDir, nsDirs)
	for _, nsInfo := range namespaceInfo {
		//log.Infof("nsInfo : %v", nsInfo)
		fmt.Println("Namespace : ", nsInfo.Name)
		for name, numOfReplica := range nsInfo.Deployments {
			fmt.Printf("   Deployment : %s  Replicas: %d\n", name, numOfReplica)
		}
		fmt.Println("")
	}
	if generate {
		for _, nsInfo := range namespaceInfo {
			//log.Infof("nsInfo : %v", nsInfo)
			for name, numOfReplica := range nsInfo.Deployments {
				fmt.Printf("       go run main.go --ports 8000 --upstream-uris http://httpbin.default:8000 --protocol http --namespace-labels istio-injection=enabled --namespace %s --deployment %s  --deployment-replicas %d\n", nsInfo.Name, name, numOfReplica)
			}
			fmt.Println("")
		}
	}
	if err != nil {
		log.Errorf("Failed to scan sub-directories, error: %w", err)
		return err
	}

	return nil
}
