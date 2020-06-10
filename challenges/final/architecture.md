Architecture Document
=====================
The way the architecture is supposed to work and the flow is going to be the following:

First thing to start is the main.go, this package will import the scheduler, controller and api.

In the main we start the controller first, where we will have the setup for workers and use the surveyor protocol
to setup the mangos sockets, each worker received from the socket (api) will be stored in a list of type Worker 
so it can later be utilized.

After setting up the controller we call the scheduler who will create the jobs for the workers and assign them 
with grpc (proto) (we use the same say hello function)

After that we get the API up and running, the API will have functions for login, logout, status, uploadImage, workerStatus, workload
the login will be authenticated and secure using the Gin framework (we will use it for the whole API really)

With all of this up and running we just need the worker to start the job filtering the images provided. this last
part wasn't finished but the code to do this is almost complete.

Hope the code is understandable, I tried to comment each section.

If it's possible to take into account the code instead of just the result I would very much appreciate it, I don't expect
a full grade, but maybe something not super bad...? hahahaha thank you either way!


