package logrusr

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"errors"
	"github.com/go-logr/logr"
	"github.com/sirupsen/logrus"
	logrus_test "github.com/sirupsen/logrus/hooks/test"
)

var _ = Describe("Logger", func() {
	var (
		logger logr.Logger
		l      *logrus.Logger
		hook   *logrus_test.Hook
	)

	BeforeEach(func() {
		l, hook = logrus_test.NewNullLogger()
		logger = New("foo", *l)
	})

	Describe("Logger Calls", func() {
		Context("For Info", func() {
			It("when standard", func() {
				logger.Info("test log", "hello", "world")

				Expect(hook.LastEntry().Message).To(Equal("test log"))
				Expect(hook.LastEntry().Level).To(Equal(logrus.InfoLevel))

				r, ok := hook.LastEntry().Data["request"].(*logrus.Fields)
				Expect(ok).To(BeTrue())

				Expect((*r)["name"]).To(Equal("foo"))
				Expect((*r)["kvs"]).To(Equal(map[string]interface{}{"hello": "world"}))
			})
			It("when named", func() {
				namedLogger := logger.WithName("bar")
				namedLogger.Info("test log", "hello", "world")

				Expect(hook.LastEntry().Message).To(Equal("test log"))
				Expect(hook.LastEntry().Level).To(Equal(logrus.InfoLevel))

				r, ok := hook.LastEntry().Data["request"].(*logrus.Fields)
				Expect(ok).To(BeTrue())

				Expect((*r)["name"]).To(Equal("foo.bar"))
				Expect((*r)["kvs"]).To(Equal(map[string]interface{}{"hello": "world"}))
			})
			It("when has values", func() {
				valuesLogger := logger.WithValues("goodbye", "crazy world")
				valuesLogger.Info("test log", "hello", "world")

				Expect(hook.LastEntry().Message).To(Equal("test log"))
				Expect(hook.LastEntry().Level).To(Equal(logrus.InfoLevel))

				r, ok := hook.LastEntry().Data["request"].(*logrus.Fields)
				Expect(ok).To(BeTrue())

				Expect((*r)["name"]).To(Equal("foo"))
				Expect((*r)["kvs"]).To(Equal(map[string]interface{}{"goodbye": "crazy world", "hello": "world"}))
			})
		})

		Context("For Err", func() {
			It("when standard", func() {
				err := errors.New("BOOM SUCKA!")

				logger.Error(err, "test error log", "hello", "world")

				Expect(hook.LastEntry().Message).To(Equal("test error log"))
				Expect(hook.LastEntry().Level).To(Equal(logrus.ErrorLevel))

				r, ok := hook.LastEntry().Data["request"].(*logrus.Fields)
				Expect(ok).To(BeTrue())

				Expect((*r)["name"]).To(Equal("foo"))
				Expect((*r)["kvs"]).To(Equal(map[string]interface{}{"hello": "world"}))
			})
		})

		Context("For Supressing", func() {
			It("when not verbose", func() {
				hook.Reset()
				logger.V(1).Info("test verbose log", "hello", "crazy world")

				Expect(hook.LastEntry()).To(BeNil())
			})
			It("when verbose", func() {
				SetVerbosity(1)

				vLogger := logger.V(1)
				vLogger.Info("test verbose log", "hello", "crazy world")

				Expect(hook.LastEntry().Message).To(Equal("test verbose log"))
				Expect(hook.LastEntry().Level).To(Equal(logrus.InfoLevel))

				r, ok := hook.LastEntry().Data["request"].(*logrus.Fields)
				Expect(ok).To(BeTrue())

				Expect((*r)["name"]).To(Equal("foo"))
				Expect((*r)["kvs"]).To(Equal(map[string]interface{}{"v": 1, "hello": "crazy world"}))
			})
			It("when limited", func() {
				hook.Reset()
				SetVerbosity(1)
				LimitToLoggers("bar")

				vLogger := logger.V(1)
				vLogger.Info("test verbose log", "hello", "crazy world")

				Expect(hook.LastEntry()).To(BeNil())
			})
		})
	})
})
