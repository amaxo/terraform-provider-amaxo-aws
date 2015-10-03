package amaxoaws

import (
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"net/url"
	"os"

	"github.com/hashicorp/terraform/terraform"
	"github.com/hashicorp/terraform/helper/schema"
	"github.com/hashicorp/terraform/builtin/providers/aws"
)

func Provider() terraform.ResourceProvider {
	provider := aws.Provider().(*schema.Provider)
	instance := provider.ResourcesMap["aws_s3_bucket_object"]
	instance.Create = wrapCreate(instance.Create)
	instance.Update = wrapCreate(instance.Update)
	return provider
}

// If the 'source' attribute is parseable as a URL, use http.Get to
// retrieve the source to a temporary file.
// Patch the original source and call create, then replace the source before returning
func wrapCreate(create func (*schema.ResourceData, interface{}) error) func (*schema.ResourceData, interface{}) error {
	wrapped := func (d *schema.ResourceData, meta interface{}) error {
		originalSource := d.Get("source").(string)

		fileUrl, err := url.Parse(originalSource)

		if err == nil && fileUrl.Scheme != "" {
			resp, err := http.Get(originalSource)
			if err != nil {
				return fmt.Errorf("Error getting S3 bucket object source (%s): %s", originalSource, err)
			}

			defer resp.Body.Close()

			temp, err := ioutil.TempFile(os.TempDir(), "download_")
			if err != nil {
				return fmt.Errorf("Error creating tempfile to download %s: %s", originalSource, err)
			}
			defer os.Remove(temp.Name())

			_, err = io.Copy(temp, resp.Body)
			if err != nil {
				return fmt.Errorf("Error copying download %s to tempfile: %s", originalSource, temp.Name(), err)
			}

			d.Set("source", temp.Name())
		}

		result := create(d, meta)

		d.Set("source", originalSource)

		return result
	}
	return wrapped
}
