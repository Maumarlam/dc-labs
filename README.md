dc-labs
=======

This repository stands as a template git repository for the **Distributed Computing** class's challenges and laboratories.

For instructions on how to submit your solutions, read [classify.md](./classify.md)

The project doesn't seem to be working because of trouble getting the correct mango import, I get the following error:

.\main.go:113:15: cannot assign "nanomsg.org/go/mangos/v2".Socket to sock (type "github.com/nanomsg/mangos".Socket) in multiple assignment:
        "nanomsg.org/go/mangos/v2".Socket does not implement "github.com/nanomsg/mangos".Socket (wrong type for Info method)
                have Info() "nanomsg.org/go/mangos/v2".ProtocolInfo
                want Info() "github.com/nanomsg/mangos".ProtocolInfo

Even though I only import the one with github... when checking my mod file it used to say that mangos was
not compatible...

I didn't finish the CUDA implementation but I added some filters from a repo I found in a tutorial about image filtering

Hope at least the conections make sense between all the different packages, they did work for me a couple times
before running into the issue stated above, it happened while I was moving and organizing files


Hope it makes sense! 
