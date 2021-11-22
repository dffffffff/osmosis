
# Install Osmosis

This guide will explain how to install the osmosid binary into your system.


On Ubuntu start by updating your system:
```bash
sudo apt update
sudo apt upgrade
```

## Install build requirements

Install make and gcc.
```bash
sudo apt install git build-essential ufw curl jq snapd --yes
```

Install go:

```bash
wget -q -O - https://git.io/vQhTU | bash -s -- --version 1.17.2
```

After installed, open new terminal to properly load go

## Install Osmosis Binary

Clone the osmosis repo, checkout and install v4.2.0 (note, if you are making a testnet node, skip this step):

```bash
cd $HOME
git clone https://github.com/osmosis-labs/osmosis
cd osmosis
git checkout v4.2.0
make install
```
 
