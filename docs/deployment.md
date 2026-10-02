### **Use Case — Deployment of Narwhal**

**Prerequisites**

*   Proxmox VE environment
    
*   Debian 13 LXC container
    
*   Network connectivity between the Proxmox node and container
    
*   Narwhal deployment files
    

#### **1\. Create the LXC container**

Create a new LXC container in Proxmox and configure:

*   Hostname
    
*   Root password
    
*   Disk size
    
*   CPU cores
    
*   Memory
    
*   Static IPv4 address
    

All other container settings can remain at their default values.

The container uses **Debian 13**.

#### **2\. Access the container**

From the Proxmox host, access the newly created container:

`pct enter`

Install the required dependencies. For example:

`apt install tmux -y`

Additional dependencies should be installed here as required by the game server being deployed.

#### **3\. Create the gameserver user**

Create the user responsible for running the game server and Narwhal.

Narwhal deployments use the following user:

`gameserver`

The game server and Narwhal should be installed within this user's home directory.

Optionally, copy an SSH key to the container if direct SSH access to the container is desired.

The SSH key is not required for Narwhal itself. However, the API requires SSH authentication to the Proxmox node, so configuring key-based authentication prevents the API from requiring a password for each interaction.

#### **4\. Deploy Narwhal**

Narwhal can be deployed using the provided `deploy.ps1` script.

Before running the script, configure:

`$deployIP = ""`

Set this value to the IP address of the target container.

#### **5\. Deploy the game server**

Install the desired game server in the gameserver user's home directory.

Narwhal should be located alongside the game server and have the necessary permissions to manage it.

#### **6\. Configure Narwhal**

Create the configuration file:

`/home/gameserver/mgr.env`

Configure the file according to the requirements of the deployed game server.

#### **7\. Test the deployment**

At this point, the deployment should be ready.

Run Narwhal manually:

`./Narwhal start`

Verify that the game server starts successfully and that Narwhal reports the expected status.

**Expected Result:**Narwhal is installed and configured inside the gameserver user's home directory and can successfully start, stop, restart, and retrieve the status of the deployed game server.