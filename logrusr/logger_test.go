package logrusr

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	// . "github.com/activeshadow/logr/logrusr"

	"errors"
	"github.com/go-logr/logr"
	"github.com/sirupsen/logrus"
	logrus_test "github.com/sirupsen/logrus/hooks/test"
	"time"
)

var _ = Describe("Logger", func() {
	var (
		logger logr.Logger
		l      *logrus.Logger
		hook   *logrus_test.Hook
		mock   time.Time
	)

	BeforeEach(func() {
		mock, _ = time.Parse("2006-01-02", "2015-12-15")
		l, hook = logrus_test.NewNullLogger()

		logger = New("foo", *l)
		logger.(*LogrusLogr).clock = &clock{mock: mock}
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
				namedLogger.(*LogrusLogr).clock = &clock{mock: mock}
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
				valuesLogger.(*LogrusLogr).clock = &clock{mock: mock}
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
				vLogger.(*LogrusLogr).clock = &clock{mock: mock}
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
				vLogger.(*LogrusInfoLogr).clock = &clock{mock: mock}
				vLogger.Info("test verbose log", "hello", "crazy world")

				Expect(hook.LastEntry()).To(BeNil())
			})
		})
	})
})
