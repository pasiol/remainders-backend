from diagrams import Cluster, Diagram
from diagrams.k8s.compute import Pod, Cronjob, StatefulSet
from diagrams.k8s.network import Service, Ingress
from diagrams.k8s.storage import PV, PVC, StorageClass
from diagrams.onprem.compute import Server
from diagrams.onprem.database import MongoDB
from diagrams.custom import Custom
from urllib.request import urlretrieve

certmanager_icon = "logo.png"
urlretrieve(
    "https://raw.githubusercontent.com/cert-manager/cert-manager/master/logo/logo.png",
    certmanager_icon,
)

with Diagram("Remainders", show=True, outformat="png"):

    mongo = Pod("MongoDB")
    frontend = Service("frontend")
    backend = Service("Backend")
    ingress = Ingress("domain.com")
    ingress - Custom("Cert Manager", certmanager_icon) - ingress
    ingress - frontend
    system_x = Server("System X")
    getter = Cronjob("Remainder getter")
    mailer = Cronjob("Remainder mailer")
    system_x >> getter >> mongo >> mailer

    with Cluster("Frontend deployment"):
        (frontend - [Pod("frontend"), Pod("frontend"), Pod("frontend")])

    (frontend - backend - mongo)
