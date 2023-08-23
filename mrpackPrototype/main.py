import json
import os
import pathlib
import hashlib
import requests

from typing import List
from enum import Enum

import xml.etree.ElementTree as Et

# First we need to be able to build out the actual data from the pack, focusing on the json aspect first

global WorkingDir


class AbsolutePathException(Exception):
    pass


class InvalidFileHashException(Exception):
    pass


class VersionNotFoundException(Exception):
    pass


class InstallerNotImplementedException(Exception):
    pass


class Hashes(dict):
    def get_sha1(self) -> str:
        return self.get("sha1")

    def get_sha512(self) -> str:
        return self.get("sha512")


class EnvEnum(Enum):
    REQUIRED = "required"
    UNSUPPORTED = "unsupported"
    OPTIONAL = "optional"


class Env(dict):
    def get_client(self) -> EnvEnum:
        return self.get("client")

    def get_server(self) -> EnvEnum:
        return self.get("server")


class File(dict):
    # noinspection PyGlobalUndefined
    def download(self):
        global response

        path = self.get_path()
        if path.is_absolute():
            raise AbsolutePathException(f"Mod {path.name} uses an Absolute path, and cannot be installed.")

        path.parent.mkdir(parents=True, exist_ok=True)
        for url in self.get_downloads():
            try:
                response = requests.get(url)
                response.raise_for_status()
                break

            except requests.exceptions.RequestException as e:
                print(f'Error occurred while downloading from {url}: {e}')

        if self.hash(response.content):
            with open(path, "wb") as f:
                f.write(response.content)

    def hash(self, b: bytes) -> bool:
        blob_hash = hashlib.sha512(b, usedforsecurity=True).hexdigest()
        if blob_hash == self.get_hashes().get_sha512():
            return True

        raise InvalidFileHashException("The expected has and the one retrieved didn't match")

    def get_path(self) -> pathlib.Path:
        return pathlib.Path(self.get("path"))

    # if we add functionality to the hash type, we need to retr hash
    def get_hashes(self) -> Hashes:
        return Hashes(self.get("hashes"))

    # same as above
    def get_env(self) -> Env:
        return self.get("env")

    def get_downloads(self) -> str:
        return self.get("downloads")

    def get_file_size(self) -> int:
        return self.get("fileSize")


class LoaderDependencies(dict):
    def __init__(self, *args, **kwargs):
        super().__init__(*args, **kwargs)

        if "forge" in self:
            self._loader = "forge"
        elif "neoforge" in self:
            self._loader = "neoforge"
        elif "quilt-loader" in self:
            self._loader = "quilt-loader"
        elif "fabric-loader" in self:
            self._loader = "fabric-loader"
        else:
            self._loader = None  # Default value if no valid loader found

    def get_loader_version(self) -> str:
        return self.get(self._loader)

    def get_loader(self) -> str:
        return self._loader

    def get_minecraft(self):
        return self.get("minecraft")


class MrPack(dict):
    def download_deps(self):
        for file in self.get_files():
            file.download()

    def get_format_version(self):
        return self.get("formatVersion")

    def get_game(self):
        return self.get("game")

    def get_version_id(self):
        return self.get("versionId")

    def get_name(self):
        return self.get("name")

    def get_summary(self):
        return self.get("summary")

    def get_files(self) -> List[File]:
        return [File(file_dict) for file_dict in self.get("files")]

    def get_loader_dependencies(self) -> LoaderDependencies:
        return LoaderDependencies(self.get("dependencies"))


class BaseLoaderInstaller(LoaderDependencies):
    def __init__(self, pack: MrPack, *args, **kwargs):
        super().__init__(*args, **kwargs)
        self.pack = pack

    def fetch_jar(self):
        raise NotImplementedError

    def install(self):
        self.fetch_jar()
        self.pack.download_deps()


class ForgeInstaller(BaseLoaderInstaller):
    def fetch_jar(self):
        # Todo(jaegyu) refine the forge XML metadata set
        # Check if the ver exists.
        loader_deps = self.pack.get_loader_dependencies()

        api = "https://maven.minecraftforge.net/net/minecraftforge/forge"
        r = requests.get(f"{api}/maven-metadata.xml")
        root = Et.fromstring(r.content)

        vers_elems = root.find("versioning/versions")
        if vers_elems is not None:
            for vers_elem in vers_elems.findall("version"):
                if loader_deps.get_loader_version() in vers_elem.text:
                    version = f"{loader_deps.get_minecraft()}-{loader_deps.get_loader_version()}"
                    query = f"{api}/{version}/forge-{version}-installer.jar"

                    r = requests.get(query)

                    with open(f"forge-{version}-installer.jar", "wb") as f:
                        f.write(r.content)

                    return

            raise VersionNotFoundException("Couldn't find the Forge/MC version described in .mrpack/dependencies")


class FabricInstaller(BaseLoaderInstaller):
    def fetch_jar(self):
        api = "https://meta.fabricmc.net"
        r = requests.get(f"{api}/v2/versions/installer")
        installers = json.loads(r.content)

        install_ver = installers[0]["version"]
        mc_ver = self.get_minecraft()
        load_ver = self.get_loader_version()

        r = requests.get(f"{api}/v2/versions/loader/{mc_ver}/{load_ver}/{install_ver}/server/jar")
        with open("server.jar", "wb") as jar:
            jar.write(r.content)


class InstallerFactory:
    def create_installer(self, pack: MrPack) -> BaseLoaderInstaller:
        loader_deps = pack.get_loader_dependencies()
        match loader_deps.get_loader():
            case "forge": return ForgeInstaller(pack)
            case "fabric-loader": return FabricInstaller(pack)

            case _:
                raise InstallerNotImplementedException(f"I haven't implemented {loader_deps.get_loader()} yet!")


def load_json(path: str) -> MrPack:
    with open(path, "r") as f:
        return MrPack(json.load(f))


def main():
    pack = load_json("mrpack/modrinth.index.json")

    factory = InstallerFactory()
    installer = factory.create_installer(pack)
    installer.install()


if __name__ == "__main__":

    WorkingDir = pathlib.Path(os.getcwd())
    main()
