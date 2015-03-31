require 'pg'
require 'uri'

class Database

  def getUri
    URI.parse(ENV['DATABASE_URL'] || 'postgres://127.0.0.1/tulospalvelu')
  end

  def createConnection
    uri = getUri
    PG.connect(uri.hostname, uri.port, nil, nil, uri.path[1..-1], uri.user, uri.password)
  end
end