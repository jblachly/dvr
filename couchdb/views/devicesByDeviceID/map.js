function(doc) {
	if (doc.type === 'device') {
		//emit(doc._id,1);
		emit(doc.DeviceID, doc.host)
	}
}
