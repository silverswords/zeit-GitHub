package github

import (
	"context"
	"net/http"

	cloudapi "github.com/silverswords/clouds/openapi/github"
	util "github.com/silverswords/clouds/pkgs/http"
	cloudpkgs "github.com/silverswords/clouds/pkgs/http/context"
	"golang.org/x/oauth2"
)

// PullsReviewsEdit updates the review summary on the specified pull request.
func PullsReviewsEdit(w http.ResponseWriter, r *http.Request) {
	var (
		github struct {
			Owner    string `json:"owner"     zeit:"required"`
			Repo     string `json:"repo"      zeit:"required"`
			Number   int    `json:"number"    zeit:"required"`
			ReviewID int64  `json:"review_id" zeit:"required"`
			Body     string `json:"body"      zeit:"required"`
		}
	)

	c := cloudpkgs.NewContext(w, r)
	err := c.ShouldBind(&github)
	if err != nil {
		c.WriteJSON(http.StatusBadRequest, cloudpkgs.H{"status": http.StatusBadRequest})
		return
	}

	err = util.Validate(&github)
	if err != nil {
		c.WriteJSON(http.StatusPreconditionRequired, cloudpkgs.H{"status": http.StatusPreconditionRequired})
		return
	}

	token := c.Request.Header
	t := token.Get("Authorization")

	ctx := context.Background()
	ts := oauth2.StaticTokenSource(
		&oauth2.Token{AccessToken: t},
	)

	tc := oauth2.NewClient(ctx, ts)
	client := cloudapi.NewAPIClient(tc)

	pull, _, err := client.Client.PullRequests.UpdateReview(ctx, github.Owner, github.Repo, github.Number, github.ReviewID, github.Body)
	if err != nil {
		c.WriteJSON(http.StatusRequestTimeout, cloudpkgs.H{"status": http.StatusRequestTimeout})
		return
	}

	c.WriteJSON(http.StatusOK, cloudpkgs.H{"status": http.StatusOK, "pull_request_review": pull})
}