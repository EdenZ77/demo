/*
Copyright 2025.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package controller

import (
	"context"
	"fmt"

	"k8s.io/apimachinery/pkg/runtime"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/log"

	demov1 "mashibing.com/demo/app/api/v1"
)

// AppReconciler reconciles a App object
type AppReconciler struct {
	client.Client
	Scheme *runtime.Scheme
}

// +kubebuilder:rbac:groups=demo.mashibing.com,resources=apps,verbs=get;list;watch;create;update;patch;delete
// +kubebuilder:rbac:groups=demo.mashibing.com,resources=apps/status,verbs=get;update;patch
// +kubebuilder:rbac:groups=demo.mashibing.com,resources=apps/finalizers,verbs=update

// Reconcile is part of the main kubernetes reconciliation loop which aims to
// move the current state of the cluster closer to the desired state.
// TODO(user): Modify the Reconcile function to compare the state specified by
// the App object against the actual cluster state, and then
// perform operations to make the cluster state reflect the state specified by
// the user.
//
// For more details, check Reconcile and its Result here:
// - https://pkg.go.dev/sigs.k8s.io/controller-runtime@v0.22.1/pkg/reconcile
func (r *AppReconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
	logger := log.FromContext(ctx)
	logger.Info("Start reconcile")

	// 实现你的业务逻辑
	// 1. 获取资源对象
	app := new(demov1.App)
	if err := r.Client.Get(ctx, req.NamespacedName, app); err != nil {
		logger.Error(err, " Get resource")
		return ctrl.Result{}, client.IgnoreNotFound(err)
	}

	// 2. 处理数据
	// 2.1 获取对应的字段
	action := app.Spec.Action
	object := app.Spec.Object

	// 2.2 拼接数据
	result := fmt.Sprintf("%s,%s", action, object)

	// 3. 创建结果
	appCopy := app.DeepCopy()
	appCopy.Status.Result = result

	// 4. 更新 status
	if err := r.Client.Status().Update(ctx, appCopy); err != nil {
		return ctrl.Result{}, err
	}

	logger.Info("End reconcile")

	return ctrl.Result{}, nil
}

// SetupWithManager sets up the controller with the Manager.
func (r *AppReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&demov1.App{}).
		Named("app").
		Complete(r)
}
