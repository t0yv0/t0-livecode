{
  description = "A flake for t0-livecode, including a Nixos service module";

  inputs = {
    nixpkgs.url = github:NixOS/nixpkgs/nixos-23.05;
  };

  outputs = { self, nixpkgs }:
    let
      ver = "0.0.1";
      package = { system }:
        let
          pkgs = import nixpkgs { system = system; };
        in
          pkgs.buildGo120Module rec {
            name = "t0-livecode-${ver}";
            version = "${ver}";
            src = ./.;
            vendorSha256 = null;
          };
    in
      {
        packages.x86_64-linux.default = package {
          system = "x86_64-linux";
        };
      };
}
