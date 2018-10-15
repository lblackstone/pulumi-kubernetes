// Copyright 2016-2018, Pulumi Corporation.  All rights reserved.

import * as k8s from "@pulumi/kubernetes";

//
// Somewhat convoluted example that tests that we can use `Output<T>` in a Chart's
// values/transforms/etc.
//

const cm = new k8s.core.v1.ConfigMap("nginx-config", {
    metadata: { name: "nginx-config" },
    data: { serviceType: "ClusterIP" }
});

const nginx = new k8s.helm.v2.Chart("simple-nginx", {
    repo: "stable",
    chart: "nginx-lego",
    version: "0.3.1",
    values: {
        // Blank out resource requests so it can run on minikube without scheduling problems. Set
        // Service type to ClusterIP so it can be run on minikube.
        nginx: {
            resources: null,
            service: { type: cm.data.apply(data => data.serviceType) }
        },
        default: { resources: null },
        lego: { resources: null }
    },
    transformations: [
        (o: any) => {
            // `creationTimestamp is a computed value. Test that the annotation
            if (o.metadata.annotations === undefined) {
                o.metadata.annotations = {};
            }
            o.metadata.annotations["cmcreation"] = cm.metadata.apply(m => m.creationTimestamp);
        }
    ]
});
