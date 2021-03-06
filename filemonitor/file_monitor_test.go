package filemonitor_test

import (
	//"fmt"
	"log"
	"math/rand"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"testing"
	"time"

	"gitlab.engr.illinois.edu/sp-box/boxsync/filemonitor"
)

func TestFileWatchWriteFile(t *testing.T) {
	t.Logf("Testing a simple file watch:\n")

	os.Mkdir("testing_tmp", 0777)
	defer os.RemoveAll("testing_tmp/")

	returnChannel := make(chan string)
	fileWatcher := filemonitor.NewWatcher(
		func(f *filemonitor.FileWatchEvent) {
			t.Logf("%s\n", f.FilePath)
			returnChannel <- f.FilePath
		})
	defer fileWatcher.Close()

	fileWatcher.AddAll("testing_tmp")
	os.Create("testing_tmp/test")

	file := <-returnChannel
	if strings.Compare(file, "testing_tmp/test") != 0 {
		t.Fail()
	}

	writeByteToFile("testing_tmp/test")
	file = <-returnChannel
	if strings.Compare(file, "testing_tmp/test") != 0 {
		t.Fail()
	}
}

func TestFileWatchRemoveFile(t *testing.T) {
}

func TestFileWatchRenameFile(t *testing.T) {
	t.Logf("Testing a simple file rename:\n")

	os.Mkdir("testing_tmp", 0777)
	defer os.RemoveAll("testing_tmp/")

	returnChannel := make(chan string)
	fileWatcher := filemonitor.NewWatcher(
		func(f *filemonitor.FileWatchEvent) {
			t.Logf("%s\n", f.FilePath)
			returnChannel <- f.FilePath
		})
	defer fileWatcher.Close()

	fileWatcher.AddAll("testing_tmp")
	os.Create("testing_tmp/test")

	file := <-returnChannel
	if strings.Compare(file, "testing_tmp/test") != 0 {
		t.Fail()
	}

	os.Rename("testing_tmp/test", "testing_tmp/test2")
	rename := <-returnChannel
	if strings.Compare(rename, "testing_tmp/test2") != 0 {
		t.Fail()
	}

	os.Mkdir("testing_tmp/test_dir", 0777)
	dir := <-returnChannel
	if strings.Compare(dir, "testing_tmp/test_dir") != 0 {
		t.Fail()
	}
	os.Create("testing_tmp/test_dir/test")
	file = <-returnChannel
	if strings.Compare(file, "testing_tmp/test_dir/test") != 0 {
		t.Fail()
	}
	os.Rename("testing_tmp/test_dir", "testing_tmp/dir_2")
	dir_rename := <-returnChannel
	if strings.Compare(dir_rename, "testing_tmp/dir_2") != 0 {
		t.Fail()
	}
	os.Remove("testing_tmp/dir_2/test")
	file = <-returnChannel
	if strings.Compare(file, "testing_tmp/dir_2/test") != 0 {
		t.Fail()
	}
}

func TestDirectoryWatchCreateFile(t *testing.T) {
}

func TestDirectoryWatchCreateDirectory(t *testing.T) {
	t.Logf("Testing a simple directory watch:\n")

	os.Mkdir("testing_tmp", 0777)
	defer os.RemoveAll("testing_tmp/")

	returnChannel := make(chan string)
	fileWatcher := filemonitor.NewWatcher(
		func(f *filemonitor.FileWatchEvent) {
			t.Logf("%s\n", f.FilePath)
			returnChannel <- f.FilePath
		})
	defer fileWatcher.Close()

	fileWatcher.AddAll("testing_tmp")
	os.Mkdir("testing_tmp/level2", 0777)

	dir := <-returnChannel
	if strings.Compare(dir, "testing_tmp/level2") != 0 {
		t.Fail()
	}

	os.Create("testing_tmp/level2/test")
	file := <-returnChannel
	if strings.Compare(file, "testing_tmp/level2/test") != 0 {
		t.Fail()
	}
}

func TestDirectoryWatchRandom(t *testing.T) {
	t.Logf("Testing a random series of ops:\n")

	events := buildRandomDirectorySequence(4000, "testing_tmp")
	ready := make(chan bool)
	defer close(ready)

	os.Mkdir("testing_tmp", 0777)
	defer os.RemoveAll("testing_tmp/")

	i := 0
	fileWatcher := filemonitor.NewWatcher(
		func(f *filemonitor.FileWatchEvent) {
			switch etype := f.Type; etype {
			case filemonitor.EvTypeCreate:
				t.Logf("Create: ")
			case filemonitor.EvTypeWrite:
				t.Logf("Write: ")
			case filemonitor.EvTypeRemove:
				t.Logf("Remove: ")
			case filemonitor.EvTypeRename:
				t.Logf("Rename: ")
			case filemonitor.EvTypeChmod:
				t.Logf("Chmod: ")
			}
			t.Logf("%s\n", f.FilePath)
			if events[i].event == EvTypeRename {
				if strings.Compare(filepath.Base(events[i].newName), filepath.Base(f.FilePath)) == 0 {
					ready <- true
					i++
				}
			} else {
				if strings.Compare(filepath.Base(events[i].filename), filepath.Base(f.FilePath)) == 0 {
					ready <- true
					i++
				}
			}
		})

	defer fileWatcher.Close()
	fileWatcher.AddAll("testing_tmp")

	err := doEventSequence(events, ready)
	if err != nil {
		log.Fatal(err)
	}
}

type DirectoryOperation struct {
	event    EventType
	fsType   FileSystemType
	filename string
	newName  string
}

type DirectoryTree struct {
	dirName     string
	dirs        map[string]*DirectoryTree
	files       map[string]bool
	numChildren int
}

type FileSystemType int

const (
	FSTypeFile FileSystemType = iota
	FSTypeDir
)

type EventType int

const (
	EvTypeCreate EventType = iota
	EvTypeWrite
	EvTypeRemove
	EvTypeRename
	EvTypeChmod
)

func writeByteToFile(file string) error {
	fp, err := os.OpenFile(file, os.O_WRONLY, 0666)
	if err != nil {
		return err
	}
	b := []byte{76, 79, 76}
	_, err = fp.Write(b)
	//t.Logf("wrote one byte to: %s\n", file)
	return err
}

func doEventSequence(ops []DirectoryOperation, ready chan bool) error {
	for _, op := range ops {
		//t.Logf("event type:%d, filetype:%d, filename:%s, newname:%s \n", op.event, op.fsType, op.filename, op.newName)
		var err error
		switch op.event {
		case EvTypeCreate:
			switch op.fsType {
			case FSTypeFile:
				_, err = os.Create(op.filename)
			case FSTypeDir:
				err = os.Mkdir(op.filename, 0777)
			}
		case EvTypeWrite:
			err = writeByteToFile(op.filename)
		case EvTypeRemove:
			switch op.fsType {
			case FSTypeFile:
				err = os.RemoveAll(op.filename)
			case FSTypeDir:
				err = os.RemoveAll(op.filename)
			}
		case EvTypeRename:
			err = os.Rename(op.filename, op.newName)
		}
		if err != nil {
			return err
		}
		_, ok := <-ready
		if !ok {
			return nil
		}
	}
	return nil
}

func walkTreeHelper(directory string, dTree *DirectoryTree, numSteps int) (int, string, FileSystemType) {
	if numSteps == 0 {
		return 0, directory + dTree.dirName + "/", FSTypeDir
	}

	for _, v := range dTree.dirs {
		numSteps--
		numStepsTmp, dir, fsType := walkTreeHelper(directory+dTree.dirName+"/", v, numSteps)
		numSteps = numStepsTmp
		if numSteps == 0 {
			return 0, dir, fsType
		}
	}

	for k := range dTree.files {
		numSteps--
		if numSteps == 0 {
			return 0, directory + dTree.dirName + "/" + k, FSTypeFile
		}
	}
	return numSteps, "", 0
}

func walkTree(root *DirectoryTree, numSteps int) (string, FileSystemType) {
	numStepsTmp, file, fileType := walkTreeHelper("", root, numSteps)
	if numStepsTmp > 0 {
	}
	return file, fileType
}

func dirOp(root *DirectoryTree, dest string, op EventType, fsType FileSystemType, newName string) {
	steps := strings.Split(dest, "/")
	var adder int = 0
	switch op {
	case EvTypeCreate:
		adder = 1
	case EvTypeRemove:
		if fsType == FSTypeFile {
			adder = -1
		} else if fsType == FSTypeDir {
			tmpSteps := steps[1:]
			rootTemp := root
			for len(tmpSteps) != 1 {
				rootTemp = rootTemp.dirs[tmpSteps[0]]
				tmpSteps = tmpSteps[1:]
			}
			adder = -1*rootTemp.numChildren - 1
		}
	}

	steps = steps[1:]
	root.numChildren += adder
	for len(steps) != 1 && steps[1] != "" {
		root = root.dirs[steps[0]]
		steps = steps[1:]
		root.numChildren += adder
	}

	switch op {
	case EvTypeCreate:
		if fsType == FSTypeFile {
			root.files[steps[0]] = true
		} else if fsType == FSTypeDir {
			root.dirs[steps[0]] = &DirectoryTree{steps[0], make(map[string]*DirectoryTree), make(map[string]bool), 0}
		}

	case EvTypeWrite:

	case EvTypeRemove:

		if fsType == FSTypeFile {
			delete(root.files, steps[0])
		} else if fsType == FSTypeDir {
			delete(root.dirs, steps[0])
		}

	case EvTypeRename:
		if fsType == FSTypeFile {
			delete(root.files, steps[0])
			root.files[newName] = true
		} else if fsType == FSTypeDir {
			root.dirs[newName] = root.dirs[steps[0]]
			root.dirs[newName].dirName = newName
			delete(root.dirs, steps[0])
		}
	}
}

func buildRandomDirectorySequence(numOps int, rootDir string) []DirectoryOperation {
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	root := &DirectoryTree{rootDir, make(map[string]*DirectoryTree), make(map[string]bool), 0}
	nameCounter := 0
	rootDir += "/"

	var events []DirectoryOperation = make([]DirectoryOperation, 0, numOps)

	for ops := 0; ops < numOps; ops++ {
		location, fileType := walkTree(root, r.Intn(root.numChildren+1))
		var choice EventType
		switch fileType {
		case FSTypeFile:
			i := r.Intn(3)
			switch i {
			case 0:
				choice = EvTypeWrite
			case 1:
				choice = EvTypeRemove
			case 2:
				choice = EvTypeRename
			}
		case FSTypeDir:
			if location == rootDir {
				choice = EvTypeCreate
				break
			}
			i := r.Intn(3)
			switch i {
			case 0:
				choice = EvTypeCreate
			case 1:
				choice = EvTypeRemove
			case 2:
				choice = EvTypeRename
			}
		}

		switch choice {
		case EvTypeCreate:
			createType := r.Intn(2)
			if createType == 0 {
				name := location + "file" + strconv.Itoa(nameCounter)
				events = append(events, DirectoryOperation{choice, FSTypeFile, name, ""})
				dirOp(root, name, EvTypeCreate, FSTypeFile, "")
			} else if createType == 1 {
				name := location + "dir" + strconv.Itoa(nameCounter)
				events = append(events, DirectoryOperation{choice, FSTypeDir, name, ""})
				dirOp(root, name, EvTypeCreate, FSTypeDir, "")
			}
			nameCounter++
		case EvTypeWrite:
			events = append(events, DirectoryOperation{choice, fileType, location, ""})
		case EvTypeRemove:
			events = append(events, DirectoryOperation{choice, fileType, location, ""})
			dirOp(root, location, EvTypeRemove, fileType, "")
		case EvTypeRename:
			if fileType == FSTypeFile {
				name := "file" + strconv.Itoa(nameCounter)
				dirOp(root, location, EvTypeRename, fileType, name)
				steps := strings.Split(location, "/")
				steps[len(steps)-1] = name
				name = strings.Join(steps, "/")
				events = append(events, DirectoryOperation{choice, fileType, location, name})
			} else if fileType == FSTypeDir {
				name := "dir" + strconv.Itoa(nameCounter)
				dirOp(root, location, EvTypeRename, fileType, name)
				steps := strings.Split(location, "/")
				steps[len(steps)-2] = name
				name = strings.Join(steps, "/")

				events = append(events, DirectoryOperation{choice, fileType, location, name})
			}
			nameCounter++
		}
	}
	return events
}
