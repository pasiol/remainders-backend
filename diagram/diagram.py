from diagrams import Cluster, Diagram
from diagrams.k8s.compute import Pod, StatefulSet
from diagrams.k8s.network import Service, Ingress
from diagrams.k8s.storage import PV, PVC, StorageClass
from diagrams.onprem.database import MongoDB
from diagrams.custom import Custom
from urllib.request import urlretrieve

certmanager_icon = "logo.png"
urlretrieve(
    "https://raw.githubusercontent.com/cert-manager/cert-manager/master/logo/logo.png",
    certmanager_icon,
)

with Diagram("Remainders", show=False):
    mongo = MongoDB("Db")
    frontend = Service("frontend")
    backend = Service("Backend")
    ingress = Ingress("domain.com")
    ingress << Custom("Cert Manager", certmanager_icon) << ingress
    ingress << frontend

    (
        frontend
        >> [Pod("frontend1"), Pod("frontend2"), Pod("frontend3")]
        >> backend
        >> mongo
    )
