package diagnostics

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"go.opencensus.io/stats/view"

	"github.com/dapr/dapr/pkg/config"
)

const (
	componentName = "test"
)

func componentsMetrics() *componentMetrics {
	c := newComponentMetrics()
	_ = c.Init("test", "default", config.LoadDefaultConfiguration().GetMetricsSpec().GetLatencyDistribution(log))

	return c
}

func TestPubSub(t *testing.T) {
	t.Run("record drop by app or sidecar", func(t *testing.T) {
		c := componentsMetrics()

		c.PubsubIngressEvent(t.Context(), componentName, "drop", "success", "A", 1)
		c.PubsubIngressEvent(t.Context(), componentName, "drop", "drop", "A", 1)

		viewData, _ := view.RetrieveData("component/pubsub_ingress/count")
		v := view.Find("component/pubsub_ingress/count")

		allTagsPresent(t, v, viewData[0].Tags)

		assert.Len(t, viewData, 2)
		assert.Equal(t, int64(1), viewData[0].Data.(*view.CountData).Value)
		assert.Equal(t, int64(1), viewData[1].Data.(*view.CountData).Value)
	})

	t.Run("record ingress count", func(t *testing.T) {
		c := componentsMetrics()

		c.PubsubIngressEvent(t.Context(), componentName, "retry", "retry", "A", 0)

		viewData, _ := view.RetrieveData("component/pubsub_ingress/count")
		v := view.Find("component/pubsub_ingress/count")

		allTagsPresent(t, v, viewData[0].Tags)
	})

	t.Run("record ingress latency", func(t *testing.T) {
		c := componentsMetrics()

		c.PubsubIngressEvent(t.Context(), componentName, "retry", "", "A", 1)

		viewData, _ := view.RetrieveData("component/pubsub_ingress/latencies")
		v := view.Find("component/pubsub_ingress/latencies")

		allTagsPresent(t, v, viewData[0].Tags)

		assert.InEpsilon(t, 1, viewData[0].Data.(*view.DistributionData).Min, 0)
	})

	t.Run("record egress latency", func(t *testing.T) {
		c := componentsMetrics()

		c.PubsubEgressEvent(t.Context(), componentName, "A", true, 1)

		viewData, _ := view.RetrieveData("component/pubsub_egress/latencies")
		v := view.Find("component/pubsub_egress/latencies")

		allTagsPresent(t, v, viewData[0].Tags)

		assert.InEpsilon(t, 1, viewData[0].Data.(*view.DistributionData).Min, 0)
	})
}

func TestBindings(t *testing.T) {
	t.Run("record input binding count", func(t *testing.T) {
		c := componentsMetrics()

		c.InputBindingEvent(t.Context(), componentName, false, 0)

		viewData, _ := view.RetrieveData("component/input_binding/count")
		v := view.Find("component/input_binding/count")

		allTagsPresent(t, v, viewData[0].Tags)
	})

	t.Run("record input binding latency", func(t *testing.T) {
		c := componentsMetrics()

		c.InputBindingEvent(t.Context(), componentName, false, 1)

		viewData, _ := view.RetrieveData("component/input_binding/latencies")
		v := view.Find("component/input_binding/count")

		allTagsPresent(t, v, viewData[0].Tags)

		assert.InEpsilon(t, 1, viewData[0].Data.(*view.DistributionData).Min, 0)
	})

	t.Run("record output binding count", func(t *testing.T) {
		c := componentsMetrics()

		c.OutputBindingEvent(t.Context(), componentName, "set", false, 0)

		viewData, _ := view.RetrieveData("component/output_binding/count")
		v := view.Find("component/input_binding/count")

		allTagsPresent(t, v, viewData[0].Tags)
	})

	t.Run("record output binding latency", func(t *testing.T) {
		c := componentsMetrics()

		c.OutputBindingEvent(t.Context(), componentName, "set", false, 1)

		viewData, _ := view.RetrieveData("component/output_binding/latencies")
		v := view.Find("component/output_binding/latencies")

		allTagsPresent(t, v, viewData[0].Tags)

		assert.InEpsilon(t, 1, viewData[0].Data.(*view.DistributionData).Min, 0)
	})
}

func TestState(t *testing.T) {
	t.Run("record state count", func(t *testing.T) {
		c := componentsMetrics()

		c.StateInvoked(t.Context(), componentName, "get", false, 0)

		viewData, _ := view.RetrieveData("component/state/count")
		v := view.Find("component/state/count")

		allTagsPresent(t, v, viewData[0].Tags)
	})

	t.Run("record state latency", func(t *testing.T) {
		c := componentsMetrics()

		c.StateInvoked(t.Context(), componentName, "get", false, 1)

		viewData, _ := view.RetrieveData("component/state/latencies")
		v := view.Find("component/state/latencies")

		allTagsPresent(t, v, viewData[0].Tags)
		assert.InEpsilon(t, 1, viewData[0].Data.(*view.DistributionData).Min, 0)
	})
}

func TestConfiguration(t *testing.T) {
	t.Run("record configuration count", func(t *testing.T) {
		c := componentsMetrics()

		c.ConfigurationInvoked(t.Context(), componentName, "get", false, 0)

		viewData, _ := view.RetrieveData("component/configuration/count")
		v := view.Find("component/configuration/count")

		allTagsPresent(t, v, viewData[0].Tags)
	})

	t.Run("record configuration latency", func(t *testing.T) {
		c := componentsMetrics()

		c.ConfigurationInvoked(t.Context(), componentName, "get", false, 1)

		viewData, _ := view.RetrieveData("component/configuration/latencies")
		v := view.Find("component/configuration/latencies")

		allTagsPresent(t, v, viewData[0].Tags)

		assert.InEpsilon(t, 1, viewData[0].Data.(*view.DistributionData).Min, 0)
	})
}

func TestSecrets(t *testing.T) {
	t.Run("record secret count", func(t *testing.T) {
		c := componentsMetrics()

		c.SecretInvoked(t.Context(), componentName, "get", false, 0)

		viewData, _ := view.RetrieveData("component/secret/count")
		v := view.Find("component/secret/count")

		allTagsPresent(t, v, viewData[0].Tags)
	})

	t.Run("record secret latency", func(t *testing.T) {
		c := componentsMetrics()

		c.SecretInvoked(t.Context(), componentName, "get", false, 1)

		viewData, _ := view.RetrieveData("component/secret/latencies")
		v := view.Find("component/secret/latencies")

		allTagsPresent(t, v, viewData[0].Tags)

		assert.InEpsilon(t, 1, viewData[0].Data.(*view.DistributionData).Min, 0)
	})
}

func TestConversation(t *testing.T) {
	t.Run("record conversation count", func(t *testing.T) {
		c := componentsMetrics()

		c.ConversationInvoked(t.Context(), componentName, false, 0)

		viewData, _ := view.RetrieveData("component/conversation/count")
		v := view.Find("component/conversation/count")

		allTagsPresent(t, v, viewData[0].Tags)
	})

	t.Run("record conversation latency", func(t *testing.T) {
		c := componentsMetrics()

		c.ConversationInvoked(t.Context(), componentName, false, 1)

		viewData, _ := view.RetrieveData("component/conversation/latencies")
		v := view.Find("component/conversation/latencies")

		allTagsPresent(t, v, viewData[0].Tags)

		assert.InEpsilon(t, 1, viewData[0].Data.(*view.DistributionData).Min, 0)
	})
}

func TestComponentMetricsInit(t *testing.T) {
	c := componentsMetrics()
	assert.True(t, c.enabled)
	assert.Equal(t, "test", c.appID)
	assert.Equal(t, "default", c.namespace)
}

func TestElapsedSince(t *testing.T) {
	start := time.Now()
	time.Sleep(time.Second)

	elapsed := ElapsedSince(start)
	assert.GreaterOrEqual(t, elapsed, float64(1000))
}
