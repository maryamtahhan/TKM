package webhooks

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/go-logr/logr"
	tkmv1 "github.com/redhat-et/TKM/api/v1alpha1"

	"github.com/google/go-containerregistry/pkg/name"
	"github.com/google/go-containerregistry/pkg/v1/remote"

	"github.com/sigstore/cosign/v2/pkg/cosign"
	csigRemote "github.com/sigstore/cosign/v2/pkg/oci/remote"

	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/webhook/admission"
)

type TKMCacheMutator struct {
	Client  client.Client
	Log     logr.Logger
	Decoder admission.Decoder
}

func (m *TKMCacheMutator) Handle(ctx context.Context, req admission.Request) admission.Response {
	m.Log.Info("handling mutation", "kind", req.Kind.Kind)

	switch req.Kind.Kind {
	case "TKMCache":
		var cache tkmv1.TKMCache
		if err := m.Decoder.Decode(req, &cache); err != nil {
			return admission.Errored(http.StatusBadRequest, err)
		}

		if cache.Spec.ResolvedDigest != "" {
			return admission.Allowed("digest already present")
		}

		digest, err := resolveImageDigest(cache.Spec.Image)
		if err != nil {
			m.Log.Error(err, "digest resolution failed")
			return admission.Denied(fmt.Sprintf("could not resolve image digest: %v", err))
		}

		if err := verifyImageSignature(ctx, cache.Spec.Image); err != nil {
			m.Log.Error(err, "image signature verification failed")
			return admission.Denied(fmt.Sprintf("signature verification failed: %v", err))
		}

		cache.Spec.ResolvedDigest = digest

		marshaled, err := json.Marshal(&cache)
		if err != nil {
			return admission.Errored(http.StatusInternalServerError, err)
		}

		return admission.PatchResponseFromRaw(req.Object.Raw, marshaled)

	case "TKMCacheCluster":
		var cluster tkmv1.TKMCacheCluster
		if err := m.Decoder.Decode(req, &cluster); err != nil {
			return admission.Errored(http.StatusBadRequest, err)
		}

		if cluster.Spec.ResolvedDigest != "" {
			return admission.Allowed("digest already present")
		}

		digest, err := resolveImageDigest(cluster.Spec.Image)
		if err != nil {
			m.Log.Error(err, "digest resolution failed")
			return admission.Denied(fmt.Sprintf("could not resolve image digest: %v", err))
		}

		if err := verifyImageSignature(ctx, cluster.Spec.Image); err != nil {
			m.Log.Error(err, "image signature verification failed")
			return admission.Denied(fmt.Sprintf("signature verification failed: %v", err))
		}

		cluster.Spec.ResolvedDigest = digest

		marshaled, err := json.Marshal(&cluster)
		if err != nil {
			return admission.Errored(http.StatusInternalServerError, err)
		}

		return admission.PatchResponseFromRaw(req.Object.Raw, marshaled)

	default:
		return admission.Allowed("unknown kind, skipping mutation")
	}
}

func resolveImageDigest(imageRef string) (string, error) {
	ref, err := name.ParseReference(imageRef)
	if err != nil {
		return "", fmt.Errorf("invalid image reference: %w", err)
	}
	desc, err := remote.Get(ref)
	if err != nil {
		return "", fmt.Errorf("failed to get image manifest: %w", err)
	}
	return desc.Descriptor.Digest.String(), nil
}

func verifyImageSignature(ctx context.Context, imageRef string) error {
	ref, err := name.ParseReference(imageRef)
	if err != nil {
		return fmt.Errorf("invalid image reference: %w", err)
	}

	co := &cosign.CheckOpts{
		ClaimVerifier:      cosign.SimpleClaimVerifier,
		RegistryClientOpts: []csigRemote.Option{},
		// keyless OIDC: Fulcio and Rekor will be auto-discovered via TUF
	}

	sigs, verified, err := cosign.VerifyImageSignatures(ctx, ref, co)
	if err != nil {
		return fmt.Errorf("cosign verification failed: %w", err)
	}
	if !verified || len(sigs) == 0 {
		return fmt.Errorf("no valid signatures found")
	}
	return nil
}

func SetupWebhook(mgr ctrl.Manager) error {
	decoder := admission.NewDecoder(mgr.GetScheme())

	hook := &TKMCacheMutator{
		Client:  mgr.GetClient(),
		Log:     ctrl.Log.WithName("webhook").WithName("TKMCache"),
		Decoder: decoder,
	}

	mgr.GetWebhookServer().Register("/mutate-tkmcache", &admission.Webhook{
		Handler: hook,
	})

	return nil
}
