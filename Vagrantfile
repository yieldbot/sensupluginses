# -*- mode: ruby -*-
# vi: set ft=ruby :

# version 0.0.1

# Vagrantfile API/syntax version. Don't touch unless you know what you're doing!
VAGRANTFILE_API_VERSION = '2'

$script = <<SCRIPT
SRCROOT="/opt/go"
SRCPATH="/opt/gopath"
# Get the ARCH
ARCH=`uname -m | sed 's|i686|386|' | sed 's|x86_64|amd64|'`
# Install Prereq Packages
sudo apt-get update
sudo apt-get install -y build-essential curl git-core libpcre3-dev mercurial pkg-config zip expect
# Install Go
cd /tmp
wget -q https://storage.googleapis.com/golang/go1.5.1.linux-${ARCH}.tar.gz
tar -xf go1.5.1.linux-${ARCH}.tar.gz
sudo mv go $SRCROOT
sudo chmod 775 $SRCROOT
sudo chown vagrant:vagrant $SRCROOT
# Setup the GOPATH; even though the shared folder spec gives the working
# directory the right user/group, we need to set it properly on the
# parent path to allow subsequent "go get" commands to work.
sudo mkdir -p $SRCPATH
sudo chown -R vagrant:vagrant $SRCPATH 2>/dev/null || true
# ^^ silencing errors here because we expect this to fail for the shared folder
cat <<EOF >/tmp/gopath.sh
export GOPATH="$SRCPATH"
export GOROOT="$SRCROOT"
export PATH="$SRCROOT/bin:$SRCPATH/bin:\$PATH"
export GO15VENDOREXPERIMENT=1
EOF
sudo mv /tmp/gopath.sh /etc/profile.d/gopath.sh
sudo chmod 0755 /etc/profile.d/gopath.sh
source /etc/profile.d/gopath.sh
touch /home/vagrant/.ssh/config chown vagrant:vagrant /home/vagrant/.ssh/config
chmod 600 /home/vagrant/.ssh/config
go get github.com/axw/gocov/gocov
go get -u github.com/golang/lint/golint
go get github.com/mitchellh/gox
cat << 'EOF' >> /home/vagrant/.ssh/config
StrictHostKeyChecking no
EOF
#exit 0
SCRIPT

Vagrant.configure(VAGRANTFILE_API_VERSION) do |config|
  config.vm.box = 'ubuntu/trusty64'
  config.vm.hostname = 'sensu-yieldbot-handler-slack'

  config.vm.provision 'shell', inline: $script, privileged: false
  config.vm.synced_folder '.', '/opt/gopath/src/github.com/yieldbot/sensu-yieldbot-handler-slack'
  config.ssh.forward_agent = true

  config.vm.provider :virtualbox do |vb|
    vb.name = config.vm.hostname
    vb.customize ['modifyvm', :id, '--memory', '2048']
    vb.customize ['modifyvm', :id, '--cpuexecutioncap', '50']
    vb.customize ['modifyvm', :id, '--cpus', '2']
    vb.customize ['modifyvm', :id, '--ioapic', 'on']
  end
end
