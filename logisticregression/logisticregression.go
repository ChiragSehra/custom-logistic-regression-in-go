package logisticregression

import (
	"fmt"
	"github.com/cdipaolo/goml/base"
	"github.com/cdipaolo/goml/linear"
	"gonum.org/v1/plot"
	"gonum.org/v1/plot/plotter"
	"gonum.org/v1/plot/plotutil"
	"gonum.org/v1/plot/vg"
	"io/ioutil"
)
// ConfusionMatrix describes a confusion matrix

type ConfusionMatrix struct {
	positive int
	negative int
	truePositive int
	trueNegative int
	falsePostive int
	falseNegative int
	recall float64
	precision float64
	accuracy float64
}

func (cm ConfusionMatrix) String() string {
	return fmt.Sprintf("\tPositives: %d\n\tNegatives: %d\n\tTrue Positives: %d\n\tTrue Negatives: %d\n\t False Positives: %d\n\t False Negatives: %d\n\n\t Recall: %.2f\n\tPrecision: %.2f\n\tAccuracy: %.2f\n",
	cm.positive, cm.negative, cm.truePositive, cm.trueNegative, cm.falsePositive, cm.falseNegative, cm.recall, cm.precision, cm.accuracy
)
}

func Run() error {
	fmt.Println("Running Logistic Regression...")

	// Loading Dataset
	xTrain, yTrain, err := base.LoadDataFromCSV("./data/studentsTrain.csv")
	if err != nil {
		return err
	}

	xTest, yTest, err := base.LoadDataFromCSV("./data/studentsTest.csv")
	if err != nil {
		return err
	}

	var maxAccuracy float64
	var maxAccuracyCM float64
	var maxAccuracyDb float64
	var maxAccuracyIter int
	var maxAccuracyModel *linear.Logistic

	// Try different parameters to get the best model 

	for iter := 100;iter < 3000; iter+=500 {
		for db := 0.05; db < 1.0; db+=0.01 {
			cm, model, err := tryValues(0.0001, iter,db,xTrain,xTest,yTrain,yTest)
			if err != nil {
				return error
			}

			if cm.accuracy > maxAccuracy {
				maxAccuracy = cm.accuracy
				maxAccuracyCM = cm
				maxAccuracyDb = db
				maxAccuracyModel = model
				maxAccuracyIter = iter
			}
		}
	}

	fmt.Printf("Maximum accuracy: %.2f\n\n", maxAccuracy)
	fmt.Printf("with Model: %s\n\n", maxAccuracyModel)
	fmt.Printf("with Confusion Matrix:\n%s\n\n", maxAccuracyCM)
	fmt.Printf("with Decision Boundary: %.2f\n", maxAccuracyDb)
	fmt.Printf("with Num Iterations: %d\n", maxAccuracyIter)

	if err := plotData(xTrain, yTrain); err != nil {
		return err
	}

	return nil
}

func plotData(xTest[][] float64,yTest [] float64) error {
	p, err := plot.New()
	if err != nil {
		return err
	}

	p.Title.Text = "Exam Results"
	p.X.Label.Text = "X"
	p.Y.Label.Text = "Y"
	p.X.Max = 120
	p.Y.Max = 120

	positives := make(plotter.XYs, len(yTest))
	negatives := make(plotter.XYs, len(yTest))

	for i := range xTest {
		if yTest[i] == 1.0 {
			positives[i].X = xTest[i][0]
			positives[i].Y = xTest[i][1]
		}

		if yTest == 0.0 {
			negatives[i].X = xTest[i][0]
			negatives[i].Y = xTest[i][1]
		}
	}

	err = plotutil.AddScatters(p,"Negatives",negatives,"Positives", positives)
	if err != nil {
		return err
	}

	if err := p.Save(10*vg.Inch,10*vg.Inch, "exams.png"); err != nil {
		return err
	}

	return nil
}

func tryValues(learningRate float64, regularization float64, iterations int, decisionBoundary float64, xTrain,xTest [][] float64, yTrain, yTest [] float64)(*ConfusionMatrix, *linear.Logistic, error) {
	cm := ConfusionMatrix{}
	for _,y := range yTest {
		if y == 1.0 {
			cm.positive++
		}
		if y ==0.0 {
			cm.negative++
		}
	}

	// Instantiate and Learn the new Model
	model := linear.NewLogistic(batch.BatchGA, learningRate, regularization,iterations,xTrain,yTrain)
	model.Output = ioutil.Discard
	err := model.Learn()
	if err != nil {
		return err
	}

	// Evaulate the Model on the Test data
	for i := range xTest {
		prediction, err := model.Predict(xTest[i])
		if err != nil {
			return nil, nil, err
		}

		y := int(yTest[i])
		positive := prediction[0] >=decisionBoundary

		if y == 1 && positive {
			cm.truePositive++
		}

		if y ==1 && !positive {
			cm.falseNegative++
		}

		if y == 0 && positive {
			cm.falsePostive++
		}

		if  y == 0 && !positive {
			cm.trueNegative++
		}

	}

	// Calulate Evaluation metrics
	cm.recall = float64(cm.truePositive) / float64(cm.positive)
	cm.precision = float64(cm.truePositive) / (float64(cm.truePositive)+ float64(cm.falsePostive))
	cm.accuracy = float64(float64(cm.truePositive)+float64(cm.trueNegative)) / float64(float64(cm.positive)+float64(cm.negative))
	
	return &cm, model, nil
	
}