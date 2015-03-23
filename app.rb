set :root, File.dirname(__FILE__)
set :views, Proc.new { File.join(root, "views") }
set :public_folder, Proc.new { File.join(root, "static") }

get '/' do
  erb :index, :layout => :layout
  #"Sinatra Heroku Cedar Template - The bare minimum for a sinatra app on cedar, running rack, and using bundler."
end

get '/kisa/:id' do
  @kisa = {}
  @kisa['nimi'] = 'Kumpulan hiihtokisat'

  erb :kisa, :layout => :layout
end

get '/kilpailija' do
  erb :kilpailija_list, :layout => :layout
end