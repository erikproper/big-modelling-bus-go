/*
 *
 * Package: mbconnect
 * Layer:   3
 * Module:  raw_artefacts
 *
 * ..... ... .. .
 *
 * Creator: Henderik A. Proper (e.proper@acm.org), TU Wien, Austria
 *
 * Version of: XX.11.2025
 *
 */

package mbconnect

const (
	rawArtefactsFilePathElement = "artefacts/file"
)

/*
 *
 * Externally visible functionality
 *
 */

func (b *TModellingBusConnector) PostRawArtefact(context, format, fileName, localFilePath string) {
	topicPath := rawArtefactsFilePathElement +
		"/" + context +
		"/" + format +
		"/" + fileName

	b.postFile(topicPath, localFilePath, GetTimestamp())
}

func (b *TModellingBusConnector) ListenForRawArtefactPostings(agentID, context, format, fileName string, postingHandler func(string)) {
	topicPath := rawArtefactsFilePathElement +
		"/" + context +
		"/" + format +
		"/" + fileName
		
	b.modellingBusEventsConnector.listenForEvents(agentID, topicPath, func(message []byte) {
		localFilePath,_ := b.getLinkedFileFromRepository(message, jsonFileName) 
		postingHandler(localFilePath)
	})
}


func (b *TModellingBusConnector) GetRawArtefact(agentID, context, format, fileName, localFileName string) string {
	topicPath := rawArtefactsFilePathElement +
		"/" + context +
		"/" + format +
		"/" + fileName

	localFilePath,_ := b.getFileFromPosting(agentID, topicPath, localFileName)
	return localFilePath
}

func (b *TModellingBusConnector) DeleteRawArtefact(context, format, fileName string) {
	topicPath := rawArtefactsFilePathElement +
		"/" + context +
		"/" + format

	b.deleteFile(topicPath, fileName)
}
