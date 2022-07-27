load 'tree.rb'

describe Node do
  context "The Node class" do
    it "should initialize correctly" do
      mock_layout = nil
      node = Node.new mock_layout
      expect(node.children).to eq []
    end

    it "should arrange itself" do
      mock_layout = double("Layout")
      node = Node.new mock_layout
      allow(mock_layout).to receive(:arrange).with(node, Position.new(1, 2, 3, 4)).and_return("arrangement")

      actual = node.arrange(Position.new(1, 2, 3, 4))

      expect(actual).to eq "arrangement"
    end

    it "should find its next sibling/node/term" do
      # root
      # - child1
      #   - child11
      #   - child12
      # - child2
      #   - child21
      root = Node.new(nil)
      child1 = Node.new(nil)
      child11 = Node.new(nil)
      child12 = Node.new(nil)
      child2 = Node.new(nil)
      child21 = Node.new(nil)

      root.add_child(child1)
      child1.add_child(child11)
      child1.add_child(child12)
      root.add_child(child2)
      child2.add_child(child21)

      main = Main.new

      # All nodes
      expect(main.all_nodes(root)).to eq [root, child1, child11, child12, child2, child21]

      # Next node
      [[root, child1],
       [child1, child11],
       [child11, child12],
       ].each do |test|
        node, expected = test
        expect(main.next_node(root, node)).to equal expected
      end

      # Next node with a block
      actual = main.next_node(root, root) {|n| n.equal? child21}
      expect(actual).to equal child21

      # Previous node
      [[child1, root],
       [child11, child1],
       ].each do |test|
        node, expected = test
        expect(main.previous_node(root, node)).to equal expected
      end

      # Previous node with a block
      actual = main.previous_node(root, child21) {|n| n.equal? child2}
      expect(actual).to equal child2
    end

    it "should let layout handle key" do
      mock_layout = double("Layout")
      node = Node.new mock_layout
      allow(mock_layout).to receive(:handle_key).with("a").and_return(true)

      expect(node.handle_key "a").to eq true
    end

    it "should let parent handle key if layout can't" do
      mock_layout = double("Layout")
      node = Node.new mock_layout
      node.parent = double("Node")
      allow(mock_layout).to receive(:handle_key).with("a").and_return(false)
      allow(node.parent).to receive(:handle_key).with("a").and_return(true)

      expect(node.handle_key "a").to eq true
    end
  end
end

describe StackLayout do
  context "The VerticalStackLayout class" do
    it "should layout correctly" do
      child1 = double("Child 1")
      allow(child1).to receive(:is_popup).and_return(false)
      allow(child1).to receive(:attribute).with("fixed-size").and_return(nil)
      allow(child1).to receive(:arrange).with(Position.new(0, 0, 80, 11)).and_return([1, 2])

      child2 = double("Child 2")
      allow(child2).to receive(:is_popup).and_return(false)
      allow(child2).to receive(:attribute).with("fixed-size").and_return(4)
      allow(child2).to receive(:arrange).with(Position.new(0, 11, 80, 4)).and_return([3, 4])

      child3 = double("Child 3")
      allow(child3).to receive(:is_popup).and_return(false)
      allow(child3).to receive(:attribute).with("fixed-size").and_return(nil)
      allow(child3).to receive(:arrange).with(Position.new(0, 15, 80, 10)).and_return([5, 6])

      parent = double("Parent")
      allow(parent).to receive(:children).and_return([child1, child2, child3])
      allow(parent).to receive(:term).and_return(nil)

      layout = VerticalStackLayout.new

      positions = layout.arrange(parent, Position.new(0, 0, 80, 25))

      expect(positions).to eq [1, 2, 3, 4, 5, 6]
    end
  end
end

describe HorizontalStackLayout do
  context "The HorizontalStackLayout class" do
    it "should layout correctly" do
      child1 = double("Child 1")
      allow(child1).to receive(:is_popup).and_return(false)
      allow(child1).to receive(:attribute).with("fixed-size").and_return(nil)
      allow(child1).to receive(:arrange).with(Position.new(0, 0, 38, 25)).and_return([1, 2])

      child2 = double("Child 2")
      allow(child2).to receive(:is_popup).and_return(false)
      allow(child2).to receive(:attribute).with("fixed-size").and_return(4)
      allow(child2).to receive(:arrange).with(Position.new(38, 0, 4, 25)).and_return([3, 4])

      child3 = double("Child 3")
      allow(child3).to receive(:is_popup).and_return(false)
      allow(child3).to receive(:attribute).with("fixed-size").and_return(nil)
      allow(child3).to receive(:arrange).with(Position.new(42, 0, 38, 25)).and_return([5, 6])

      parent = double("Parent")
      allow(parent).to receive(:is_popup).and_return(false)
      allow(parent).to receive(:children).and_return([child1, child2, child3])
      allow(parent).to receive(:term).and_return(nil)

      layout = HorizontalStackLayout.new

      positions = layout.arrange(parent, Position.new(0, 0, 80, 25))

      expect(positions).to eq [1, 2, 3, 4, 5, 6]
    end
  end
end

describe BaseLayout do
end

describe Term do
  context "The Term" do
    it "should invoke the api to create the panel" do
      client = double("Client")
      allow(client).to receive(:send).and_return(123)
      term = Term.new(client, ["/bin/bash"])
      expect(term.id).to eq 123
    end

#    it "should layout itself correctly" do
#      client = double("Client")
#      allow(client).to receive(:send).and_return(123)
#      term = Term.new(client, ["/bin/bash"])
#
#      positions = term.arrange(0, 0, 80, 25)
#
#      expected = [PanelState.new(term.panel_id, 0, 0, 80, 25)]
#      expect(positions).to eq expected
#    end
  end
end

describe Main do
  context "The main function" do
    #it "should json" do
    #root = Main.new.start nil

    #positions = root.arrange(0, 0, 80, 25)

    #json = JSON.generate(positions)
    #puts json
    #expect(json).to eq []
    #end

#    it "should yaml" do
#      root = Main.new.start nil
#
#      positions = root.arrange(0, 0, 80, 25)
#
#      yml = YAML::dump(Main.new.simplify(positions))
#      puts yml
#      expect(yml).to eq []
#    end
  end
end
