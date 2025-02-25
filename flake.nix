{
  description = "A flake for t0-livecode, including a Nixos service module";

  inputs = {
    nixpkgs.url = github:NixOS/nixpkgs/nixos-24.11;
  };

  outputs = { self, nixpkgs }: let

    ver = "0.0.1";

    p5 = "https://cdnjs.cloudflare.com/ajax/libs/p5.js/1.7.0/";
    codeMirror = "https://cdnjs.cloudflare.com/ajax/libs/codemirror/6.65.7";

    htmx = builtins.fetchurl {
      name = "htmx";
      url = "https://unpkg.com/htmx.org@2.0.4";
      sha256 = "0ixlixv36rrfzj97g2w0q6jxbg0x1rswgvvd2vrpjm13r2jxs2g2";
    };

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

    assets = pkgs: pkgs.stdenv.mkDerivation {
      name = "t0-assets-${ver}";
      src = ./.;
        buildPhase = ''
        mkdir -p $out
        mkdir -p $out/www
        cp ${cmCSS}    $out/www/codemirror.min.css
        cp ${cmJS}     $out/www/codemirror.min.js
        cp ${cmModeJS} $out/www/javascript.min.js
        cp ${p5JS}     $out/www/p5.min.js
        cp ${htmx}     $out/www/htmx.min.js
      '';
    };

    fullSources = pkgs: pkgs.stdenv.mkDerivation {
      name = "t0-livecode-sources-${ver}";
      src = ./.;
      buildPhase = ''
        mkdir -p $out
        cp -r $src/*.go $src/*.mod $src/*.sum $out/
        mkdir -p $out/www
        cp ${assets pkgs}/www/* $out/www/
        cp -r $src/www/* $out/www/
      '';
    };

    app = pkgs: pkgs.buildGo123Module rec {
      name = "t0-livecode-${ver}";
      version = "${ver}";
      src = fullSources pkgs;
      vendorHash = null;
    };

    package = thing: system: let
      pkgs = import nixpkgs { system = system; };
    in thing pkgs;
  in {
    packages.x86_64-linux.default = package app "x86_64-linux";
    packages.aarch64-darwin.default = package app "aarch64-darwin";
    packages.x86_64-linux.assets = package assets "x86_64-linux";
    packages.aarch64-darwin.assets = package assets "aarch64-darwin";

    nixosModules.t0-livecode-module = { config, pkgs, ... }:
      with nixpkgs.lib; let
        cfg = config.services.t0-livecode;
        p = app pkgs;
      in {
        imports = [];
        options = {
          services.t0-livecode = {
            enable = mkEnableOption "t0-livecode";
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
            };
          };
        };
      };
  };
}
