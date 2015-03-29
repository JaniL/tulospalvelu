require '../lib/database'

class Kilpailija

  def initialize(rivi=nil)
    if rivi != nil
      sarakkeet = rivi.keys
      sarakkeet.each { |sarake| self.instance_variable_set(sarake,rivi[sarake]) }
    end
  end

  def all
    yhteys = Database.new.createConnection

    res = yhteys.exec('SELECT * FROM Kilpailija')
  end

  def find(id)
    yhteys = Database.new.createConnection

    res = yhteys.exec_params('SELECT * FROM Kilpailija WHERE id == a', [id])
  end

  def save

  end
end