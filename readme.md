#FSGrid
FSGrid is an extension package of the functional golang version of FBP [1]Grids and it implements the necessary file system structures for use in the functional fbp pattern

##Install
    
    go get github.com/influx/fsgrid
 
 Then
 
    go install github.com/influx/fsgrid

##API
Fsgrid ensures to implement as much basic FileSystem structures that can be composed to get the different functionality desired and it will continousely evolve to ensure best practices. The api is driven to provided majorly three class of operations based on Files,Directory and Control.

####Directories: 
This structs composed the Grids API and create specific grid based struct that handle directory reading and writing operations and for easier use there exists simple instance creating functions to simplify its use:

-  FSDir: This grid struct handle the reading and writing of directory by providing two channels (“Read” and “Write”) that recieves grid packets that contain as <header> file paths and meta details and in the case of writers streams of data as <body> to be written. Readers and Writers structs are built based on this struct(compose on this struct) as it encapsulates the basic operations.

    Helper Functions include:

    1.  CreateFSDir: This creates a FSDir struct  and returns a point which has included all the necessary logic for the reading and writing of directories,its logic is written in a way to allow multiple use i.e not bound to specific paths so as to allow it to be like a reader and writer request manager, writing or reading to or from directory paths within the packet `<header>`  and in the case of a writer, the Packet packets sequence is treated as subdirectories to be created in the specified path

    2.   ReadDir(file string): This helper creates a DirReader struct which enforces a one directory reading operation i.e this creates a Reader that encapsulates and on servers the read channel of a FSDir grid for read only operations.

    3.   WriteDir(file string): This helper creates a DirWriter struct which enforces a one directory writing operation and maps all its operations down to the FSDir write channel. If the given Path in the <header> does not exists, it will be created and the packets will be treated as subdirectories to be created within the path given

    Examples:
           
    -    Using FSDir direcctly:

               ```
                    file := CreateFSDir()
                    packet := grids.NewPacket()
                    packet.Set("file", "./fsgrid")

                    ev := file.Out("res")
                    ev.Receive(func(i interface{}) {
                        //do something
                    })

                    file.InSend("read", packet)
               
               ```

    -    Using ReadDir:

               ```
                    file := ReadDir(“./boom”)

                    //create an empty packet ,sort of a kicker
                    packet := grids.CreatePacket()

                    file.OrOut("res",func (p *GridPacket){
                        //do something 
                        //the body is an immute.Sequence hence you can filter,map,... etc on sequence operations
                        var seq = p.Body;

                        //helper for calling seq.Each(...)
                        seq.Offload(...)

                        //or 

                        seq.map(...)
                    })

                    file.InSend("read", packet)
               
               ```

    -    Using WriteDir:

               ```
                    file := WriteDir(“./boom”)

                    //create an empty packet ,sort of a kicker
                    packet := grids.CreatePacket()
                    //these will be treated as subdirectories to be created under “./boom”
                    packet.Push(“apps”)
                    packet.Push(“config”)

                    //lets freeze this packet so it does not allow any other push ,not necessary but its a nice locking 
                    //mechinism for those desiring

                    packet.Freeze()

                    file.OrOut("res",func (p *GridPacket){
                        //do something 
                        //the body is an immute.Sequence hence you can filter,map,... etc on sequence operations
                        var seq = p.Body;

                        //helper for calling seq.Each(...)
                        seq.Offload(...)

                        //or 

                        seq.map(...)
                    })

                    file.InSend("read", packet)
               
               ```

####Files: 
This structs composed the Grids API and create specific grid based struct that handle file reading and writing operations and for easier use there exists simple instance creating functions to simplify its use:

 - FSFile: 
   This grid struct handle the reading and writing of files by providing two channels (“Read” and “Write”) that recieves grid packets that contain as <header> file paths and meta details and in the case of writers streams of data as <body> to be written. Readers and Writers structs are built based on this struct(compose on this struct) as it encapsulates the basic file operations.

 Helper Functions include:

    1.  CreateFSFile: This creates a FSFile struct  and returns a point which has included all the necessary logic for the reading and writing of files,its logic is written in a way to allow multiple use i.e not bound to specific file paths so as to allow it to be like a file reader and writer request manager, writing or reading to or from file paths listen within the packet `<header>`

    2.   ReadFile(file string): This helper creates a FileReader struct which enforces a one file reading operation i.e this creates a Reader that encapsulates and on servers the read channel of a FSFile grid for read only operations.

    3.   WriteFile(file string): This helper creates a FileWriter struct which enforces a one file writing operation and maps all its operations down to the FSFile write channel

    Examples:
           
    -    Using FSFile direcctly:

               ```
                    file := CreateFSFile()
                    packet := grids.NewPacket()
                    packet.Set("file", "./fsgrid.go")

                    ev := file.Out("res")
                    ev.Receive(func(i interface{}) {
                        //do something
                    })

                    file.InSend("read", packet)
               
               ```

   -    Using ReadFile:

               ```
                    file := ReadFile(“./boom.tx”)

                    //create an empty packet ,sort of a kicker
                    packet := grids.CreatePacket()

                    file.OrOut("res",func (p *GridPacket){
                        //do something 
                        //the body is an immute.Sequence hence you can filter,map,... etc on sequence operations
                        var seq = p.Body;

                        //helper for calling seq.Each(...)
                        seq.Offload(...)

                        //or 

                        seq.map(...)
                    })

                    file.InSend("read", packet)
               
               ```

    -    Using WriteFile:

               ```
                    file := ReadFile(“./boom.tx”)

                    //create an empty packet ,sort of a kicker
                    packet := grids.CreatePacket()
                    packet.Push(“so who cares?”)
                    packet.Push(“i do.”)

                    //lets freeze this packet so it does not allow any other push ,not necessary but its a nice locking 
                    //mechinism for those desiring

                    packet.Freeze()

                    file.OrOut("res",func (p *GridPacket){
                        //do something 
                        //the body is an immute.Sequence hence you can filter,map,... etc on sequence operations
                        var seq = p.Body;

                        //helper for calling seq.Each(...)
                        seq.Offload(...)

                        //or 

                        seq.map(...)
                    })

                    file.InSend("read", packet)
               
               ```
