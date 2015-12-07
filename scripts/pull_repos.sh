#! /bin/bash

# This will pull all golang package dependencies.

# It will compile a unique list using `go list` and then attempt to pull them
# using `go get`. Due to limitations with this tool it will only handle https
# requests, therefore requiring a username/password combo to pull private repos.
# To get around this we can use expect to detect if this is the case and then
# mark that package to be cloned into the gopath via ssh. This will require that
# you have some method of dropping keys, in Vagrant you can use ssh-agent and
# then set `forward_agent` to true in the Vagrantfile.

PKGLIST=$(go list ./... | xargs go list -f '{{join .Deps "\n"}}' | sort -u)

for i in $PKGLIST; do
  /usr/bin/expect << EOF &> /dev/null
      spawn go get $i
      expect "Username for 'https://github.com':"
      send "\r"
EOF
      if [ $? == 0 ]; then
        if [ $(echo $i | awk -F/ '{ print $NF }') == src ]; then
          install_dir=$(echo $i | sed 's#/src##')
          install_pkg=$(echo git@$i | sed 's#/src#.git#' | sed 's#/#:#')
        else
          install_pkg=$(echo git@$i.git | sed 's#/#:#')
          install_dir=$(echo $i | sed 's#/src##')
        fi

        echo RUNNING: git clone $install_pkg $GOPATH/src/$install_dir
        git clone $install_pkg $GOPATH/src/$install_dir
      fi
done
