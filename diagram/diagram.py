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

    mongo = Pod("mongo")
    mongo_svc = Service("mongo-nodeport-svc")
    frontend_svc = Service("frontend")
    backend_svc = Service("Backend")
    ingress = Ingress("domain.com")
    ingress - Custom("Cert Manager", certmanager_icon) - ingress
    ingress - frontend_svc
    system_x = Server("System X")
    getter = Cronjob("Remainder getter")
    mailer = Cronjob("Remainder mailer")
    system_x >> getter >> mongo >> mailer
    backend = Pod("backend")
    pv = PVC("mongo-data-pv")
    pvc = PVC("mongo-data-pvc")
    backend_svc - backend - mongo_svc - mongo -pvc -pv

    with Cluster("Frontend deployment"):
        (
            frontend_svc
            - [Pod("frontend"), Pod("frontend"), Pod("frontend")]
            - backend_svc
        )

    

    mailer - Service("mailhog-smtp-svc") - Pod("mailhog") - Service("mailhog-web-svc") - Ingress("mailhog-ing")

   


