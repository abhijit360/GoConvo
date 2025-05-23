package main
import "testing"

func TestAuthAvatar(t *testing.T){
	var authAvatar AuthAvatar
	client := new(client)
	url, err := authAvatar.GetAvatarURL(client)

	if err != ErrNoAvatar {
		t.Error("AuthAvatar.GetAvatarURL should return ErrNoAvatarURL when no value present")
	}

	testURL := "http://mock-url/"
	client.userData = map[string]interface{}{"avatar_url": testURL}
	url, err = authAvatar.GetAvatarURL(client)
	if err != nil {
		t.Error("AuthAvatar.GetAvatarURL should return no error when value present")
	}

	if url != testURL {
		t.Error("AuthAvatar.GetAvatarURL should return correct URL")
	}

}

func TestGravatarAvatar(t *testing.T){
	var gravatar GravatarAvatar
	client := new(client)
	client.userData = map[string]interface{}{"email: MyEmailAddress@example.com"}
	url, err := gravatarAvatar.GetAvatarURL(client)
	if err != nil {
		t.Error("GravatarAvatar.GetAvatarURL should not return an error")
	}
	if url != "//www.gravatar.com/avatar/0bc83cb571cd1c50ba6f3e8a78ef1346" {
		t.Errorf("GravatarAvatar.GetAvatarURL wrong return %s",url)
	}
}