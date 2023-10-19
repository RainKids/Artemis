package sentinel

import (
	"fmt"

	sentinel "github.com/alibaba/sentinel-golang/api"
	"github.com/alibaba/sentinel-golang/core/circuitbreaker"
	"github.com/alibaba/sentinel-golang/core/config"
	"github.com/alibaba/sentinel-golang/core/flow"
	"github.com/alibaba/sentinel-golang/logging"
	ginSentinel "github.com/alibaba/sentinel-golang/pkg/adapters/gin"
	"github.com/alibaba/sentinel-golang/util"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

type stateChangeTestListener struct {
}

func (s *stateChangeTestListener) OnTransformToClosed(prev circuitbreaker.State, rule circuitbreaker.Rule) {
	fmt.Printf("rule.steategy: %+v, From %s to Closed, time: %d\n", rule.Strategy, prev.String(), util.CurrentTimeMillis())
}

func (s *stateChangeTestListener) OnTransformToOpen(prev circuitbreaker.State, rule circuitbreaker.Rule, snapshot interface{}) {
	fmt.Printf("rule.steategy: %+v, From %s to Open, snapshot: %.2f, time: %d\n", rule.Strategy, prev.String(), snapshot, util.CurrentTimeMillis())
}

func (s *stateChangeTestListener) OnTransformToHalfOpen(prev circuitbreaker.State, rule circuitbreaker.Rule) {

}

var resName = "GET:/api/v1/post"

func InitFlowQPS() {
	if err := sentinel.InitDefault(); err != nil {
		zap.S().Error()
	}
	_, err := flow.LoadRules([]*flow.Rule{
		{
			Resource:               resName,
			TokenCalculateStrategy: flow.Direct,
			ControlBehavior:        flow.Reject,
			Threshold:              5,
			StatIntervalInMs:       6000,
		},
	})
	if err != nil {
		zap.S().Errorf("Unexpected error: %+v", err)
		return
	}
}

func FlowQPSMiddleware() gin.HandlerFunc {
	return ginSentinel.SentinelMiddleware(
		ginSentinel.WithBlockFallback(func(ctx *gin.Context) {
			ctx.AbortWithStatusJSON(200, map[string]interface{}{
				"msg":  "too many request; the quota used up!",
				"code": 500,
			})
		}),
	)
}

func InitCircuitBreaker() {
	conf := config.NewDefaultConfig()
	// for testing, logging output to console
	conf.Sentinel.Log.Logger = logging.NewConsoleLogger()
	err := sentinel.InitWithConfig(conf)
	//if err := sentinel.InitDefault(); err != nil {
	//	zap.S().Error()
	//}
	if err != nil {
		fmt.Println("ssss", err)
		zap.S().Error()
	}
	circuitbreaker.RegisterStateChangeListeners(&stateChangeTestListener{})

	_, err = circuitbreaker.LoadRules([]*circuitbreaker.Rule{
		// Statistic time span=5s, recoveryTimeout=3s, maxErrorRatio=40%
		{
			Resource:                     resName,
			Strategy:                     circuitbreaker.ErrorRatio,
			RetryTimeoutMs:               3000,
			MinRequestAmount:             10,
			StatIntervalMs:               5000,
			StatSlidingWindowBucketCount: 10,
			Threshold:                    0.4,
		},
	})
	if err != nil {
		zap.S().Errorf("Unexpected error: %+v", err)
		return
	}
}
