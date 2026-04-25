package background

import (
	"github.com/CFF4HA/Dashboard/internal/core"
	"github.com/CFF4HA/Dashboard/internal/types"
)

type TaggingPipeline struct {
	Jobs chan *types.TaggingJob
}

func (pipeline *TaggingPipeline) RunTaggingPipeline() {
	for job := range pipeline.Jobs {
		core.Logger.Info("starting a pipeline job", "tag", job.TagRule.Tag.Name)

		var labels []types.Label
		if job.TagRule.Contains {
			core.DB.Model(&types.Label{}).Preload("Ingredients").Where("payload ~* ?", job.TagRule.Regex).Find(&labels)
		} else {
			core.DB.Model(&types.Label{}).Preload("Ingredients").Not("payload ~* ?", job.TagRule.Regex).Find(&labels)
		}

		for _, label := range labels {
			for _, ingredient := range label.Ingredients {
				// now we are going to tag this ingredient according to
				// the tag by placing it in the many2many table.
				core.DB.Model(&job.TagRule.Tag).Association("Ingredients").Append(&ingredient)
			}
		}
	}
}
