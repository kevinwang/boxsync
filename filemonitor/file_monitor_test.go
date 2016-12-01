package filemonitor_test

import (
    "testing"
    "time"
    "math/rand"
    //"path/filepath"
    "os"
    "strconv"
    "strings"
    "log"
    "fmt"
)

func TestFileWatchWriteFile(t *testing.T) {
}

func TestFileWatchRemoveFile(t *testing.T) {
}

func TestFileWatchRenameFile(t *testing.T) {
}

func TestDirectoryWatchCreateFile(t * testing.T) {
}

func TestDirectoryWatchCreateDirectory(t * testing.T) {
}

func TestDirectoryWatchRandom(t * testing.T) {
    os.Mkdir("tmp", 0777)
    events := buildRandomDirectorySequence(100, "tmp")
    err := doEventSequence(events)
    if err != nil {
        log.Fatal(err)
    }
}

type DirectoryOperation struct {
    event EventType
    fsType FileSystemType
    filename string
    newName string
}

type DirectoryTree struct {
    dirName string
    dirs map[string]*DirectoryTree
    files map[string]bool
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

func doEventSequence (ops []DirectoryOperation) error {
    for _, op := range ops {
        var err error
        switch op.event {
        case EvTypeCreate:
            switch op.fsType {
            case FSTypeFile:
                _, err = os.Create(op.filename)
            case FSTypeDir:
                fmt.Printf("%s\n", op.filename)
                err = os.Mkdir(op.filename, 0777)
            }
        case EvTypeWrite:
            fp, err := os.OpenFile(op.filename, os.O_WRONLY, 0666)
            if err != nil {
                return err
            }
            b := []byte{1}
            _, err = fp.Write(b)
        case EvTypeRemove:
            switch op.fsType {
            case FSTypeFile:
                err = os.Remove(op.filename)
            case FSTypeDir:
                err = os.RemoveAll(op.filename)
            }
        case EvTypeRename:
            err = os.Rename(op.filename, op.newName)
        }
        if err != nil {
            return err
        }
    }
    return nil
}

func walkTreeHelper (directory string, dTree *DirectoryTree, numSteps int) (int, string, FileSystemType) {
    if numSteps == 0 {
        return 0, directory + dTree.dirName + "/", FSTypeDir
    }

    for _, v := range dTree.dirs {
        numSteps--
        numStepsTmp, dir, fsType := walkTreeHelper(directory + dTree.dirName + "/", v, numSteps)
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

func walkTree (root *DirectoryTree, numSteps int) (string, FileSystemType) {
    numStepsTmp, file, fileType := walkTreeHelper("", root, numSteps)
    if numStepsTmp > 0 {
    }
    return file, fileType
}

func dirOp (root *DirectoryTree, dest string, op EventType, fsType FileSystemType, newName string) {
    fmt.Printf("dest: %s\n", dest)
    steps:= strings.Split(dest, "/")
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
            adder = -1 * rootTemp.numChildren - 1
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
    case EvTypeCreate :
        fmt.Printf("Create: %s\n", dest)
            if fsType == FSTypeFile {
                root.files[steps[0]] = true
            } else if fsType == FSTypeDir {
                root.dirs[steps[0]] = &DirectoryTree{steps[0], make(map[string]*DirectoryTree), make(map[string]bool), 0}
            }

    case EvTypeWrite :
        fmt.Printf("Write to: %s\n", dest)

    case EvTypeRemove :
        fmt.Printf("Remove: %s\n", dest)

        if fsType == FSTypeFile {
            delete(root.files, steps[0])
        } else if fsType == FSTypeDir {
            delete(root.dirs, steps[0])
        }

    case EvTypeRename :
        fmt.Printf("Rename %s to %s\n", dest, newName)
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

func buildRandomDirectorySequence(numOps int, rootDir string) []DirectoryOperation{
    r := rand.New(rand.NewSource(time.Now().UnixNano()))
    root := &DirectoryTree{rootDir, make(map[string]*DirectoryTree), make(map[string]bool), 0}
    nameCounter := 0
    rootDir+= "/"

    var events []DirectoryOperation = make([]DirectoryOperation, 0, numOps)

    for ops:= 0; ops < numOps; ops++ {
        location, fileType := walkTree(root, r.Intn(root.numChildren + 1))
        var choice EventType
        switch fileType {
        case FSTypeFile :
            i:= r.Intn(3)
            switch i {
            case 0 :
                choice = EvTypeWrite
            case 1 :
                choice = EvTypeRemove
            case 2 :
                choice = EvTypeRename
            }
        case FSTypeDir :
            if location == rootDir {
                choice = EvTypeCreate
                break
            }
            i:= r.Intn(3)
            switch i {
            case 0 :
                choice = EvTypeCreate
            case 1 :
                choice = EvTypeRemove
            case 2 :
                choice = EvTypeRename
            }
        }

        switch choice {
        case EvTypeCreate :
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
        case EvTypeWrite :
            events = append(events, DirectoryOperation{choice, fileType, location, ""})
        case EvTypeRemove :
            events = append(events, DirectoryOperation{choice, fileType, location, ""})
            dirOp(root, location, EvTypeRemove, fileType, "")
        case EvTypeRename :
            if fileType == FSTypeFile {
                name := "file" + strconv.Itoa(nameCounter)
                dirOp(root, location, EvTypeRename, fileType, name)
                events = append(events, DirectoryOperation{choice, fileType, location, name})
            } else if fileType == FSTypeDir {
                name := "dir" + strconv.Itoa(nameCounter)
                dirOp(root, location, EvTypeRename, fileType, name)
                events = append(events, DirectoryOperation{choice, fileType, location, name})
            }
            nameCounter++
        }
    }
    return events
}
