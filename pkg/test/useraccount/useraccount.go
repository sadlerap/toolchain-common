package useraccount

import (
	toolchainv1alpha1 "github.com/codeready-toolchain/api/api/v1alpha1"
	"github.com/codeready-toolchain/toolchain-common/pkg/condition"
	"github.com/codeready-toolchain/toolchain-common/pkg/test"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

type UaModifier func(ua *toolchainv1alpha1.UserAccount)

func NewUserAccountFromMur(mur *toolchainv1alpha1.MasterUserRecord, modifiers ...UaModifier) *toolchainv1alpha1.UserAccount {
	ua := &toolchainv1alpha1.UserAccount{
		ObjectMeta: metav1.ObjectMeta{
			Name:      mur.Name,
			Namespace: test.MemberOperatorNs,
			Labels: map[string]string{
				toolchainv1alpha1.TierLabelKey: mur.Spec.TierName,
			},
			Annotations: map[string]string{
				toolchainv1alpha1.UserEmailAnnotationKey: mur.Annotations[toolchainv1alpha1.MasterUserRecordEmailAnnotationKey],
			},
		},
		Spec: toolchainv1alpha1.UserAccountSpec{
			UserID:           mur.Spec.UserID,
			Disabled:         mur.Spec.Disabled,
			PropagatedClaims: mur.Spec.PropagatedClaims,
		},
	}
	Modify(ua, modifiers...)
	return ua
}

func Modify(ua *toolchainv1alpha1.UserAccount, modifiers ...UaModifier) {
	for _, modify := range modifiers {
		modify(ua)
	}
}

func StatusCondition(con toolchainv1alpha1.Condition) UaModifier {
	return func(ua *toolchainv1alpha1.UserAccount) {
		ua.Status.Conditions, _ = condition.AddOrUpdateStatusConditions(ua.Status.Conditions, con)
	}
}

func ResourceVersion(resVer string) UaModifier {
	return func(ua *toolchainv1alpha1.UserAccount) {
		ua.ResourceVersion = resVer
	}
}

// DisabledUa creates a UaModifier to change the disabled spec value
func DisabledUa(disabled bool) UaModifier {
	return func(ua *toolchainv1alpha1.UserAccount) {
		ua.Spec.Disabled = disabled
	}
}

// DeletedUa creates a UaModifier to set the deletion timestamp on the UserAccount
func DeletedUa() UaModifier {
	return func(ua *toolchainv1alpha1.UserAccount) {
		now := metav1.Now()
		ua.DeletionTimestamp = &now
	}
}
