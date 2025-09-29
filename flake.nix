{
    description = "Draftea Dev Env configration";

    inputs = {
        nixpkgs.url = github:nixos/nixpkgs?ref=nixos-25.05;
        flake-utils.url = github:numtide/flake-utils;
    };

    outputs = { self, nixpkgs, flake-utils }:
        flake-utils.lib.eachDefaultSystem(system:
            let
                pkgs = import nixpkgs {
                    inherit system;
                };
            in {
                packages = {
                    default = pkgs.mkShell {
                        name = "draftea-dev-env";

                        buildInputs = [
                            pkgs.go
                            pkgs.go-mockery
                            pkgs.mage
                            pkgs.revive
                            pkgs.pre-commit
                            pkgs.commitizen
                        ];

                        shellHook = ''
                            export GOROOT=${pkgs.go}/share/go
                        '';
                    };
                };
            }
        );
}
