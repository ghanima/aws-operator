package certs

type File struct {
	AbsolutePath string
	Data         []byte
}

type Files []File

func NewFilesCluster(cluster Cluster) Files {
	common := newFilesClusterCommon(cluster)
	master := newFilesClusterMaster(cluster)
	worker := newFilesClusterWorker(cluster)

	all := Files{}
	all = append(all, common...)
	all = append(all, master...)
	all = append(all, worker...)

	return all
}

func NewFilesClusterMaster(cluster Cluster) Files {
	common := newFilesClusterCommon(cluster)
	master := newFilesClusterMaster(cluster)

	all := Files{}
	all = append(all, common...)
	all = append(all, master...)

	return all
}

func NewFilesClusterWorker(cluster Cluster) Files {
	common := newFilesClusterCommon(cluster)
	worker := newFilesClusterWorker(cluster)

	all := Files{}
	all = append(all, common...)
	all = append(all, worker...)

	return all
}

func newFilesClusterCommon(cluster Cluster) Files {
	return Files{
		// TODO(r7vme): Only used by Calico and should be removed
		// when Calico will migrate to the ones below.
		{
			AbsolutePath: "/etc/kubernetes/ssl/etcd/client-ca.pem",
			Data:         cluster.EtcdServer.CA,
		},
		{
			AbsolutePath: "/etc/kubernetes/ssl/etcd/client-crt.pem",
			Data:         cluster.EtcdServer.Crt,
		},
		{
			AbsolutePath: "/etc/kubernetes/ssl/etcd/client-key.pem",
			Data:         cluster.EtcdServer.Key,
		},
		// Calico Etcd client.
		{
			AbsolutePath: "/etc/kubernetes/ssl/calico/etcd-ca",
			Data:         cluster.CalicoEtcdClient.CA,
		},
		{
			AbsolutePath: "/etc/kubernetes/ssl/calico/etcd-cert",
			Data:         cluster.CalicoEtcdClient.Crt,
		},
		{
			AbsolutePath: "/etc/kubernetes/ssl/calico/etcd-key",
			Data:         cluster.CalicoEtcdClient.Key,
		},
	}
}

func newFilesClusterMaster(cluster Cluster) Files {
	return Files{
		// Kubernetes API server.
		{
			AbsolutePath: "/etc/kubernetes/ssl/apiserver-ca.pem",
			Data:         cluster.APIServer.CA,
		},
		{
			AbsolutePath: "/etc/kubernetes/ssl/apiserver-crt.pem",
			Data:         cluster.APIServer.Crt,
		},
		{
			AbsolutePath: "/etc/kubernetes/ssl/apiserver-key.pem",
			Data:         cluster.APIServer.Key,
		},
		// Etcd server.
		{
			AbsolutePath: "/etc/kubernetes/ssl/etcd/server-ca.pem",
			Data:         cluster.EtcdServer.CA,
		},
		{
			AbsolutePath: "/etc/kubernetes/ssl/etcd/server-crt.pem",
			Data:         cluster.EtcdServer.Crt,
		},
		{
			AbsolutePath: "/etc/kubernetes/ssl/etcd/server-key.pem",
			Data:         cluster.EtcdServer.Key,
		},
		// Service account (only key file is used).
		{
			AbsolutePath: "/etc/kubernetes/ssl/service-account-key.pem",
			Data:         cluster.ServiceAccount.Key,
		},
	}
}

func newFilesClusterWorker(cluster Cluster) Files {
	return Files{
		// Kubernetes worker.
		{
			Data:         cluster.Worker.CA,
			AbsolutePath: "/etc/kubernetes/ssl/worker-ca.pem",
		},
		{
			Data:         cluster.Worker.Crt,
			AbsolutePath: "/etc/kubernetes/ssl/worker-crt.pem",
		},
		{
			Data:         cluster.Worker.Key,
			AbsolutePath: "/etc/kubernetes/ssl/worker-key.pem",
		},
	}
}
