## Logger
Logger is to log site events, primarily the user actions in to database for data analysis.

For example, in order payment successful page, a list of recommended products are displayed there. When user clicks one product link at this page, this `click` event will be logged.

Frontend will send an API request to backend for this action.

Kafka will listen to this event and produce a message for this event. Then later, the kafka consumer will log it into database. Depending on the project preference, the database could be SQL, or No-SQL.
