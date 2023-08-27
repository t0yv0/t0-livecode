{
  description = "A flake for t0-livecode, including a Nixos service module";

  inputs = {
    nixpkgs.url = github:NixOS/nixpkgs/nixos-23.05;
  };

  outputs = { self, nixpkgs }:
    let
      ver = "0.0.1";
      p5 = "https://cdnjs.cloudflare.com/ajax/libs/p5.js/1.7.0/";
      codeMirror = "https://cdnjs.cloudflare.com/ajax/libs/codemirror/6.65.7";

      cmCSS = builtins.fetchurl {
        url = "${codeMirror}/codemirror.min.css";
        sha256 = "1i2zr35sihmfw9sp832p1bqa5m67cliwd18hwjgx4mb9mc9721qi";
      };

      cmJS = builtins.fetchurl {
        url = "${codeMirror}/codemirror.min.js";
        sha256 = "1in7flnhkaarq9qqihyc8bxz5ynvdrgfzmn4cgivj41f3v78k1j5";
      };

      cmModeJS = builtins.fetchurl {
        url = "${codeMirror}/mode/javascript/javascript.min.js";
        sha256 = "1462qdkm40ak0x96x0i1ngw1l3l92rc0acdl8lfpgv9aiya0alxk";
      };

      p5JS = builtins.fetchurl {
        url = "${p5}/p5.min.js";
        sha256 = "1fmrq3cqca4y2dvyrgzc86avjphh5rpw1mjwzx226bnfp4a8yzxv";
      };

      fullSources = pkgs:
        pkgs.stdenv.mkDerivation {
          name = "t0-livecode-sources-${ver}";
          src = ./.;
          buildPhase = ''
            mkdir -p $out
            cp -r $src/*.go $src/*.mod $src/*.sum $out/
            mkdir -p $out/www
            cp ${cmCSS} $out/www/codemirror.min.css
            cp ${cmJS} $out/www/codemirror.min.js
            cp ${cmModeJS} $out/www/javascript.min.js
            cp ${p5JS} $out/www/p5.min.js
            cp -r $src/www/* $out/www/
          '';
        };

      app = pkgs:
        pkgs.buildGo120Module rec {
          name = "t0-livecode-${ver}";
          version = "${ver}";
          src = fullSources pkgs;
          vendorSha256 = null;
        };

      package = system:
        let
          pkgs = import nixpkgs { system = system; };
        in
          app pkgs;
    in
      {
        packages.x86_64-linux.default = package "x86_64-linux";

        nixosModules.t0-livecode-module = { config, pkgs, ... }:
          with nixpkgs.lib;
          let
            cfg = config.services.t0-livecode;
            p = app pkgs;
          in
            {
              imports = [];
              options = {
                servies.t0-livecode = {
                  enable = mkEnableOption "t0-livecode";
                  environment-file = mkOption {
                    description = "Root-readable env file, can store secrets";
                    example = "/etc/t0-livecode.conf";
                  };
                };
              };
              config = mkIf cfg.enable {
                systemd.services.t0-livecode = {
                  after = [ "network-online.target" ];
                  before = [ "multi-user.target" ];
                  wantedBy = [ "multi-user.target" ];
                  serviceConfig = {
                    DynamicUser = true;
                    # App can read this value with $STATE_DIRECTORY
                    StateDirectory = "t0-livecode-state";
                    ExecStart = "${p}/bin/t0-livecode";
                    ProtectSystem = "strict";
                    ProtectHome = true;
                    EnvironmentFile = cfg.environment-file;
                  };
                };
              };
            };
      };
}
