class BaseModel

  @@sarakkeet = []

  def initialize(rivi=nil)
    @@sarakkeet.each  { |sarake| self.instance_variable_set('@' + sarake,nil) }
    if rivi != nil
      sarakkeet = rivi.keys
      sarakkeet.each { |sarake| self.instance_variable_set('@' + sarake,rivi[sarake]) }
      sarakkeet.each do |sarake|
        attr_accessor('@' + sarake)
        self.instance_variable_set('@' + sarake,rivi[sarake])

      end
    end
  end

  def all
    yhteys = Database.new.createConnection

    res = yhteys.exec('SELECT * FROM ' + yhteys.quote_ident(self.class.name))
    res.map { |rivi| self.class.new rivi }
  end

  def find(id)
    yhteys = Database.new.createConnection

    res = yhteys.exec_params('SELECT * FROM ' + self.class.name + ' WHERE id = $1', [id])
    if res[0]
      self.class.new res[0]
    else
      nil
    end
  end

  def save
    yhteys = Database.new.createConnection

    res = yhteys.exec_params('INSERT INTO ' + yhteys.quote_ident(self.class.name) + ' ($2) VALUES ($3)',[sarakemerkkijono])
    res
  end

  private

  def sarakemerkkijono
    @@sarakkeet.inject { |m,n| m + ',' + n }
  end

  def arvomerkkijono
    merkkijono = ''
    @@sarakeet.each do |sarake|
      merkkijono += self.instance_variable_get('@' + sarake).to_s + ',' if self.instance_variable_get('@' + sarake) != nil
      merkkijono += "''" + ',' if self.instance_variable_get('@' + sarake) == nil
    end
    merkkijono[0..-2]
  end
end