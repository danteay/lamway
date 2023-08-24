{
    description = "Draftea Dev Env configration";

    inputs = {
        nixpkgs.url = github:nixos/nixpkgs?ref=nixos-22.11;
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
                            pkgs.go_1_20
                            pkgs.go-migrate
                            pkgs.go-mockery
                            pkgs.mage
                            pkgs.revive
                            pkgs.python311
                            pkgs.nodejs-18_x
                            pkgs.pre-commit
                            pkgs.git
                            pkgs.awscli2
                            pkgs.commitizen
                            pkgs.fish
                        ];

                        GOROOT = "${pkgs.go_1_20}/share/go";

                        shellHook = ''
                            fish

                            # Git config
                            git config --global url.git@github.com:Drafteame.insteadOf https://github.com/Drafteame

                            # The path to this repository
                            shell_nix="''${IN_LORRI_SHELL:-$(pwd)/shell.nix}"
                            workspace_root=$(dirname "$shell_nix")
                            export WORKSPACE_ROOT="$workspace_root"

                            # We put the $GOPATH/$GOCACHE/$GOENV in $TOOLCHAIN_ROOT,
                            # and ensure that the GOPATH's bin dir is on our PATH so tools
                            # can be installed with `go install`.
                            #
                            # Any tools installed explicitly with `go install` will take precedence
                            # over versions installed by Nix due to the ordering here.

                            export TOOLCHAIN_ROOT="$workspace_root/.nix"
                            export GOCACHE="$TOOLCHAIN_ROOT/go/cache"
                            export GOENV="$TOOLCHAIN_ROOT/go/env"
                            export GOPATH="$TOOLCHAIN_ROOT/go/path"
                            export GOMODCACHE="$GOPATH/pkg/mod"
                            export GOBIN="$GOPATH/bin"
                            export GO111MODULE="on"
                            export GOSUMDB="off"
                            export CGO_ENABLED="0"

                            # Installing goimports-reviser
                            go install github.com/incu6us/goimports-reviser/v3@v3.3.1

                            # NodeJS env vars
                            # export NPM_CONFIG_PREFIX="$HOME/.npm-global";

                            # PATH
                            export PATH="$GOROOT/bin:$GOBIN:$NPM_CONFIG_PREFIX/bin:$PATH";
                        '';
                    };
                };
            }
        );
}
