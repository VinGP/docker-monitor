# -*- mode: ruby -*-
# vi: set ft=ruby :

Vagrant.configure("2") do |config|
    config.vm.box = "ubuntu/jammy64"  # Ubuntu 22.04
    config.vm.network "private_network", type: "dhcp"
    config.vm.network "forwarded_port", guest: 80, host: 8080

    config.vm.provider "virtualbox" do |vb|
      vb.memory = "2048"
      vb.cpus = 2
    end

    config.vm.provision "shell", inline: <<-SHELL
      echo "🔧 Устанавливаем Docker..."
      sudo apt-get update
      sudo apt-get install ca-certificates curl -y
      sudo install -m 0755 -d /etc/apt/keyrings
      sudo curl -fsSL https://download.docker.com/linux/ubuntu/gpg -o /etc/apt/keyrings/docker.asc
      sudo chmod a+r /etc/apt/keyrings/docker.asc

      echo "deb [arch=$(dpkg --print-architecture) signed-by=/etc/apt/keyrings/docker.asc] https://download.docker.com/linux/ubuntu \
      $(. /etc/os-release && echo "${UBUNTU_CODENAME:-$VERSION_CODENAME}") stable" | \
      sudo tee /etc/apt/sources.list.d/docker.list > /dev/null
      sudo apt-get update

      sudo apt-get install docker-ce docker-ce-cli containerd.io docker-buildx-plugin docker-compose-plugin -y

      echo "🔧 Устанавливаем make..."
      sudo apt install make -y

      echo "🚀 Запускаем Docker..."
      sudo systemctl enable --now docker
      sudo usermod -aG docker vagrant


      echo "🚀 Запускаем Docker Monitor..."
      cd /vagrant
#       export
      sudo NEXT_PUBLIC_API_URL="http://localhost:8080/api" docker compose up --build -d
    SHELL
end
