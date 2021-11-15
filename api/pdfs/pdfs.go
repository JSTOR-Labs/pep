package pdfs

import (
	"encoding/gob"
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/rs/zerolog/log"

	"github.com/JSTOR-Labs/pep/api/searchtree"
	"github.com/spf13/viper"
)

func getPDFInfo(path string) (string, string, error) {
	cwd, err := os.Getwd()
	if err != nil {
		return "", "", err
	}

	parts := strings.Split(path, string(os.PathSeparator))
	// log.Println(parts)
	doi := parts[len(parts)-2]
	filenameParts := strings.Split(parts[len(parts)-1], ".")
	if len(filenameParts) < 2 {
		return "", "", errors.New("invalid filename")
	}
	var relpath string
	if viper.GetBool("runtime.flash_drive_mode") {
		relpath = path
	} else {
		relpath, err = filepath.Rel(cwd, path)
		if err != nil {
			return "", "", err
		}
	}

	doi = fmt.Sprintf("%s/%s", doi, strings.Join(filenameParts[:len(filenameParts)-1], "."))

	return doi, relpath, nil
}

func getTopic(root, path string) (string, error) {
	// fmt.Println("root: ", root, " path: ", path)
	relPath, err := filepath.Rel(root, path)
	if err != nil {
		return "", err
	}
	// fmt.Println("rel path: ", relPath)
	parts := strings.Split(relPath, string(os.PathSeparator))

	return parts[0], nil
}

func GenerateIndex(root string) error {
	// Generate index of PDFs
	bst := new(searchtree.BinarySearchTree)
	topicSizes := make(map[string]int64)

	err := filepath.Walk(root, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if !info.IsDir() {
			doi, relpath, err := getPDFInfo(path)
			if err != nil {
				log.Printf("failed to get PDF info for %s: %v", path, err)
				return nil
			}
			bst.Insert(doi, searchtree.Item(relpath))

			// Record size
			topic, err := getTopic(root, path)
			if err != nil {
				log.Printf("failed to get PDF topic for %s: %v", path, err)
				return nil
			}
			if size, ok := topicSizes[topic]; ok {
				topicSizes[topic] = size + info.Size()
			} else {
				topicSizes[topic] = info.Size()
			}
		}

		return nil
	})
	if err != nil {
		log.Error().Err(err).Msg("failed to walk PDFs")
		return err
	}

	si, err := os.Create("pdfindex.dat")
	if err != nil {
		log.Error().Err(err).Msg("failed to create pdfindex.dat")
		return err
	}
	defer si.Close()
	enc := gob.NewEncoder(si)
	err = enc.Encode(bst)
	if err != nil {
		return err
	}

	sizesIndex, err := os.Create("pdfsizes.json")
	if err != nil {
		log.Error().Err(err).Msg("failed to create pdfsizes.json")
		return err
	}
	defer sizesIndex.Close()

	jEnc := json.NewEncoder(sizesIndex)
	err = jEnc.Encode(topicSizes)
	if err != nil {
		return err
	}

	return nil
}

func GetPDFSizes() (map[string]int64, error) {
	f, err := os.Open("pdfsizes.json")
	if err != nil {
		return nil, err
	}

	dec := json.NewDecoder(f)

	var sizes map[string]int64

	err = dec.Decode(&sizes)
	if err != nil {
		return nil, err
	}

	return sizes, nil
}

func LoadIndex() (*searchtree.BinarySearchTree, error) {
	f, err := os.Open("pdfindex.dat")
	if err != nil {
		return nil, err
	}

	dec := gob.NewDecoder(f)

	bst := new(searchtree.BinarySearchTree)

	err = dec.Decode(bst)
	if err != nil {
		return nil, err
	}

	return bst, nil
}
