require 'pg'
require 'uri'

class Database

  def getUri
    URI.parse(ENV['DATABASE_URL'] || "")
  end

  def createConnection
    uri = getUri
    PG.connect(uri.hostname, uri.port, nil, nil, uri.path[1..-1], uri.user, uri.password)
  end
end