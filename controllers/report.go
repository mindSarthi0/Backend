// / Start of Selection
package controller

import (
	"myproject/models"

	"github.com/kamva/mgm/v3"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type ReportResponse struct {
	Report   []models.Report      `json:"report"`
	AiReport []models.FinalReport `json:"aiReport"`
}

// Start Generation Here
func GetReportsByTestId(testId string) (ReportResponse, error) {
	oid, err := primitive.ObjectIDFromHex(testId)
	if err != nil {
		return ReportResponse{}, err
	}

	var reports []models.Report
	if err := mgm.Coll(&models.Report{}).SimpleFind(&reports, bson.M{"testId": oid}); err != nil {
		return ReportResponse{}, err
	}

	var finalReports []models.FinalReport
	if err := mgm.Coll(&models.FinalReport{}).SimpleFind(&finalReports, bson.M{"testId": oid}); err != nil {
		return ReportResponse{}, err
	}

	return ReportResponse{Report: reports, AiReport: finalReports}, nil
}
