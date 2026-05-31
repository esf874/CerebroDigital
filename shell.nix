let
  nixpkgs = fetchTarball "https://github.com/NixOS/nixpkgs/tarball/nixos-25.05";
  pkgs = import nixpkgs { 
  config = { 
	allowUnfree = true; 
	# anhadido porque consideraba la version de mongo insegura, necesita permiso explicito
	permittedInsecurePackages = [ 
		"mongodb-7.0.25"
	];
  };
 overlays = []; };
in pkgs.mkShell {
  name = "cerebro-digital";
  packages = with pkgs; [
    mongodb
    mongosh
    go
    air
    llama-cpp

    # start-mongo: ensure db dir and fork mongod
    (writeShellScriptBin "start-mongo" ''
      #!/usr/bin/env bash
      set -euo pipefail
      mkdir -p ./data/db
      if pgrep -x mongod >/dev/null 2>&1; then
        echo "mongod already running"
        exit 0
      fi
      echo "Starting mongod with dbpath=./data/db (logs -> ./data/mongod.log)"
      mongod --dbpath ./data/db --logpath ./data/mongod.log --fork || echo "failed to start mongod"
    '')

    # stop-mongo: stop mongod
    (writeShellScriptBin "stop-mongo" ''
      #!/usr/bin/env bash
      set -euo pipefail
      if pgrep -x mongod >/dev/null 2>&1; then
        echo "Stopping mongod"
        pkill -x mongod
      else
        echo "mongod is not running"
      fi
    '')
  ];

  shellHook = ''
    # ensure a local data dir for mongod
    mkdir -p ./data/db
  '';
}
