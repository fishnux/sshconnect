# sshconnect

Do you have multiple Linux servers and miss the option to easily select the server you want to SSH into? Would you rather avoid adding server aliases to your SSH config file?

Say hello to sshconnect!

It works by reading ~/.sshconnect.json file containing a list of your servers. Upon running you're presented to that list and you're allowed to type your choice. sshconnect will execute ssh and replace itself with the newly created ssh process.

To get started just download a pre-compiled binary (or compile it yourself) and put it in your /usr/local/bin path with execute permissions:
```
# For MacOS:
wget https://github.com/fishnux/sshconnect/releases/download/v0.1.0/sshconnect_osx -O /usr/local/bin/sshconnect && chmod +x ./usr/local/bin/sshconnect
```
```
sshconnect
```
![screenshot](https://i.imgur.com/ZVATwux.png)
(sorry for the large screenshot)
